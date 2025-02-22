package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/redhatinsights/edge-api/config"
	"github.com/redhatinsights/edge-api/pkg/db"
	"github.com/redhatinsights/edge-api/pkg/errors"
	"github.com/redhatinsights/edge-api/pkg/routes/common"

	"github.com/redhatinsights/edge-api/pkg/services"

	"github.com/golang/mock/gomock"
	"github.com/redhatinsights/edge-api/pkg/models"
	"github.com/redhatinsights/edge-api/pkg/services/mock_services"

	"github.com/redhatinsights/edge-api/pkg/dependencies"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("DeviceGroup routes", func() {
	var (
		ctrl                    *gomock.Controller
		mockDeviceGroupsService *mock_services.MockDeviceGroupsServiceInterface
		edgeAPIServices         *dependencies.EdgeAPIServices
		deviceGroupName         = "test-device-group"
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockDeviceGroupsService = mock_services.NewMockDeviceGroupsServiceInterface(ctrl)
		edgeAPIServices = &dependencies.EdgeAPIServices{
			DeviceGroupsService: mockDeviceGroupsService,
			Log:                 log.NewEntry(log.StandardLogger()),
		}
		Expect(ctrl).ToNot(BeNil())
		Expect(mockDeviceGroupsService).ToNot(BeNil())
		Expect(edgeAPIServices).ToNot(BeNil())
	})
	AfterEach(func() {
		ctrl.Finish()
	})
	Context("get all devices with filter parameters", func() {
		tt := []struct {
			name          string
			params        string
			expectedError []validationError
		}{
			{
				name:   "bad created_at date",
				params: "created_at=today",
				expectedError: []validationError{
					{Key: "created_at", Reason: `parsing time "today" as "2006-01-02": cannot parse "today" as "2006"`},
				},
			},
			{
				name:   "bad sort_by",
				params: "sort_by=test",
				expectedError: []validationError{
					{Key: "sort_by", Reason: "test is not a valid sort_by. Sort-by must be name or created_at or updated_at"},
				},
			},
		}

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		for _, te := range tt {
			req, err := http.NewRequest("GET", fmt.Sprintf("/device-groups?%s", te.params), nil)
			Expect(err).ToNot(HaveOccurred())
			w := httptest.NewRecorder()
			validateGetAllDeviceGroupsFilterParams(next).ServeHTTP(w, req)

			resp := w.Result()
			var jsonBody []validationError
			err = json.NewDecoder(resp.Body).Decode(&jsonBody)
			Expect(err).ToNot(HaveOccurred())
			for _, exErr := range te.expectedError {
				found := false
				for _, jsErr := range jsonBody {
					if jsErr.Key == exErr.Key && jsErr.Reason == exErr.Reason {
						found = true
						break
					}
				}
				Expect(found).To(BeTrue(), fmt.Sprintf("in %q: was expected to have %v but not found in %v", te.name, exErr, jsonBody))
			}
		}
	})
	Context("get DeviceGroup by id", func() {
		It("should return 200", func() {
			fakeID, _ := faker.RandomInt(1000, 2000, 1)
			fakeIDUint := uint(fakeID[0])
			req, err := http.NewRequest("GET", "/", nil)
			Expect(err).To(BeNil())

			ctx := context.WithValue(req.Context(), deviceGroupKey, &models.DeviceGroup{
				Model: models.Model{
					ID: fakeIDUint,
				},
			})
			req = req.WithContext(ctx)
			ctx = dependencies.ContextWithServices(req.Context(), &dependencies.EdgeAPIServices{})
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(GetDeviceGroupByID)
			handler.ServeHTTP(rr, req)
			// Check the status code is what we expect.
			Expect(rr.Code).To(Equal(http.StatusOK))
		})
	})
	Context("get DeviceGroup by invalid id", func() {
		It("should return 400", func() {
			req, err := http.NewRequest("GET", "/", nil)
			Expect(err).To(BeNil())

			ctx := context.WithValue(req.Context(), deviceGroupKey, "a")
			req = req.WithContext(ctx)
			ctx = dependencies.ContextWithServices(req.Context(), &dependencies.EdgeAPIServices{})
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(GetDeviceGroupByID)
			handler.ServeHTTP(rr, req)
			// Check the status code is what we expect.
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})
	})
	Context("get all devices", func() {
		req, err := http.NewRequest("GET", "/", nil)
		Expect(err).To(BeNil())
		When("all is valid", func() {
			It("should return 200", func() {
				ctx := req.Context()
				ctx = dependencies.ContextWithServices(ctx, edgeAPIServices)
				req = req.WithContext(ctx)
				rr := httptest.NewRecorder()

				// setup mock for DeviceGroupsService
				mockDeviceGroupsService.EXPECT().GetDeviceGroupsCount(gomock.Any(), gomock.Any()).Return(int64(0), nil)
				mockDeviceGroupsService.EXPECT().GetDeviceGroups(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&[]models.DeviceGroupListDetail{}, nil)

				handler := http.HandlerFunc(GetAllDeviceGroups)
				handler.ServeHTTP(rr, req)
				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusOK))
			})
		})
	})
	Context("adding devices to DeviceGroup", func() {
		account := faker.UUIDHyphenated()
		deviceGroupName := faker.Name()
		devices := []models.Device{
			{
				Name:    faker.Name(),
				UUID:    faker.UUIDHyphenated(),
				Account: account,
			},
			{
				Name:    faker.Name(),
				UUID:    faker.UUIDHyphenated(),
				Account: account,
			},
			{
				Name:    faker.Name(),
				UUID:    faker.UUIDHyphenated(),
				Account: account,
			},
		}
		deviceGroup := models.DeviceGroup{Name: deviceGroupName, Account: account, Type: models.DeviceGroupTypeDefault}
		Context("adding Devices & DeviceGroup to DB", func() {
			for _, device := range devices {
				dbResult := db.DB.Create(&device).Error
				Expect(dbResult).To(BeNil())
			}
			dbResult := db.DB.Create(&deviceGroup).Error
			Expect(dbResult).To(BeNil())
		})

		Context("get DeviceGroup from DB", func() {
			dbResult := db.DB.Where(models.DeviceGroup{Account: account}).First(&deviceGroup).Error
			Expect(dbResult).To(BeNil())
			dbResult = db.DB.Where(models.Device{Account: account}).Find(&devices).Error
			Expect(dbResult).To(BeNil())
		})
		jsonDeviceBytes, err := json.Marshal(models.DeviceGroup{Devices: devices})
		Expect(err).To(BeNil())

		url := fmt.Sprintf("/%d/devices", deviceGroup.ID)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonDeviceBytes))
		Expect(err).To(BeNil())

		When("all is valid", func() {
			It("should add devices to DeviceGroup", func() {
				ctx := req.Context()
				ctx = setContextDeviceGroup(ctx, &deviceGroup)
				ctx = dependencies.ContextWithServices(ctx, edgeAPIServices)
				req = req.WithContext(ctx)
				rr := httptest.NewRecorder()

				// setup mock for DeviceGroupsService
				mockDeviceGroupsService.EXPECT().AddDeviceGroupDevices(account, deviceGroup.ID, gomock.Any()).Return(&devices, nil)

				handler := http.HandlerFunc(AddDeviceGroupDevices)
				handler.ServeHTTP(rr, req)
				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusOK))
			})
		})
	})
	Context("create DeviceGroup", func() {
		When("all is valid", func() {
			deviceGroup := &models.DeviceGroup{
				Name:    deviceGroupName,
				Type:    models.DeviceGroupTypeDefault,
				Account: common.DefaultAccount,
			}
			jsonDeviceBytes, err := json.Marshal(deviceGroup)
			Expect(err).To(BeNil())

			url := fmt.Sprintf("/%d", deviceGroup.ID)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonDeviceBytes))
			Expect(err).To(BeNil())
			It("should create DeviceGroup", func() {
				ctx := req.Context()
				ctx = setContextDeviceGroup(ctx, deviceGroup)
				ctx = dependencies.ContextWithServices(ctx, edgeAPIServices)
				req = req.WithContext(ctx)
				rr := httptest.NewRecorder()

				// setup mock for DeviceGroupsService
				mockDeviceGroupsService.EXPECT().CreateDeviceGroup(deviceGroup).Return(deviceGroup, nil)

				handler := http.HandlerFunc(CreateDeviceGroup)
				handler.ServeHTTP(rr, req)
				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusOK))
			})
		})
		When("no account", func() {
			deviceGroup := &models.DeviceGroup{
				Name:    faker.Name(),
				Type:    models.DeviceGroupTypeDefault,
				Account: "",
			}
			jsonDeviceBytes, err := json.Marshal(deviceGroup)
			Expect(err).To(BeNil())

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonDeviceBytes))
			Expect(err).To(BeNil())
			It("should return 400", func() {
				config.Get().Auth = true // enable auth to avoid default account
				ctx := req.Context()
				ctx = setContextDeviceGroup(ctx, deviceGroup)
				ctx = dependencies.ContextWithServices(ctx, edgeAPIServices)
				req = req.WithContext(ctx)
				rr := httptest.NewRecorder()

				handler := http.HandlerFunc(CreateDeviceGroup)
				handler.ServeHTTP(rr, req)
				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusBadRequest))
				config.Get().Auth = false // disable auth
			})
		})
	})
	Context("update DeviceGroup", func() {
		deviceGroupUpdated := &models.DeviceGroup{
			Name:    deviceGroupName,
			Type:    models.DeviceGroupTypeDefault,
			Account: common.DefaultAccount,
		}
		jsonDeviceBytes, err := json.Marshal(deviceGroupUpdated)
		Expect(err).To(BeNil())

		url := fmt.Sprintf("/%d", deviceGroupUpdated.ID)
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonDeviceBytes))
		Expect(err).To(BeNil())

		When("all is valid", func() {
			It("should update DeviceGroup", func() {
				ctx := req.Context()
				ctx = setContextDeviceGroup(ctx, deviceGroupUpdated)
				ctx = dependencies.ContextWithServices(ctx, edgeAPIServices)
				req = req.WithContext(ctx)
				rr := httptest.NewRecorder()

				// setup mock for DeviceGroupsService
				mockDeviceGroupsService.EXPECT().GetDeviceGroupByID(fmt.Sprintf("%d", deviceGroupUpdated.ID)).Return(deviceGroupUpdated, nil)
				mockDeviceGroupsService.EXPECT().UpdateDeviceGroup(deviceGroupUpdated, common.DefaultAccount, fmt.Sprintf("%d", deviceGroupUpdated.ID)).Return(nil)

				handler := http.HandlerFunc(UpdateDeviceGroup)
				handler.ServeHTTP(rr, req)
				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusOK))
			})
		})
	})
	Context("delete DeviceGroup", func() {
		account := common.DefaultAccount
		deviceGroupName := faker.Name()
		devices := []models.Device{
			{
				Name:    faker.Name(),
				UUID:    faker.UUIDHyphenated(),
				Account: account,
			},
			{
				Name:    faker.Name(),
				UUID:    faker.UUIDHyphenated(),
				Account: account,
			},
		}
		deviceGroup := &models.DeviceGroup{
			Name:    deviceGroupName,
			Type:    models.DeviceGroupTypeDefault,
			Account: account,
			Devices: devices,
		}
		Context("saving DeviceGroup", func() {
			dbResult := db.DB.Create(&deviceGroup).Error
			Expect(dbResult).To(BeNil())
		})
		Context("getting DeviceGroup", func() {
			dbResult := db.DB.Where(models.DeviceGroup{Name: deviceGroupName, Account: account}).First(&deviceGroup).Error
			Expect(dbResult).To(BeNil())
		})
		When("all is valid", func() {
			url := fmt.Sprintf("/%d", deviceGroup.ID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			Expect(err).To(BeNil())

			It("should return status code 200", func() {
				ctx := req.Context()
				ctx = setContextDeviceGroup(ctx, deviceGroup)
				ctx = dependencies.ContextWithServices(ctx, edgeAPIServices)
				req = req.WithContext(ctx)
				rr := httptest.NewRecorder()

				// setup mock for DeviceGroupsService
				mockDeviceGroupsService.EXPECT().DeleteDeviceGroupByID(fmt.Sprintf("%d", deviceGroup.ID)).Return(nil)

				handler := http.HandlerFunc(DeleteDeviceGroupByID)
				handler.ServeHTTP(rr, req)
				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusOK))
			})
		})
		When("no device group in context", func() {
			url := fmt.Sprintf("/%d", deviceGroup.ID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			Expect(err).To(BeNil())

			It("should return status code 400", func() {
				ctx := req.Context()
				ctx = dependencies.ContextWithServices(ctx, edgeAPIServices)
				req = req.WithContext(ctx)
				rr := httptest.NewRecorder()

				handler := http.HandlerFunc(DeleteDeviceGroupByID)
				handler.ServeHTTP(rr, req)
				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})
		})
		When("no account", func() {
			fakeID, _ := faker.RandomInt(1000, 2000, 1)
			fakeIDUint := uint(fakeID[0])
			url := fmt.Sprintf("/%d", fakeIDUint)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			Expect(err).To(BeNil())

			It("should return status code 400", func() {
				ctx := req.Context()
				ctx = setContextDeviceGroup(ctx, &models.DeviceGroup{
					Model: models.Model{
						ID: fakeIDUint,
					},
					Account: "",
				})
				ctx = dependencies.ContextWithServices(ctx, edgeAPIServices)
				req = req.WithContext(ctx)
				rr := httptest.NewRecorder()

				// setup mock for DeviceGroupsService
				mockDeviceGroupsService.EXPECT().DeleteDeviceGroupByID(fmt.Sprint(fakeIDUint)).Return(new(services.AccountNotSet))

				handler := http.HandlerFunc(DeleteDeviceGroupByID)
				handler.ServeHTTP(rr, req)
				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})
		})
		When("no such ID", func() {
			fakeID, _ := faker.RandomInt(1000, 2000, 1)
			fakeIDUint := uint(fakeID[0])
			url := fmt.Sprintf("/%d", fakeIDUint)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			Expect(err).To(BeNil())

			It("should return status code 404", func() {
				ctx := req.Context()
				ctx = setContextDeviceGroup(ctx, &models.DeviceGroup{
					Model: models.Model{
						ID: fakeIDUint,
					},
				})
				ctx = dependencies.ContextWithServices(ctx, edgeAPIServices)
				req = req.WithContext(ctx)
				rr := httptest.NewRecorder()

				// setup mock for DeviceGroupsService
				mockDeviceGroupsService.EXPECT().DeleteDeviceGroupByID(fmt.Sprint(fakeIDUint)).Return(new(services.DeviceGroupNotFound))

				handler := http.HandlerFunc(DeleteDeviceGroupByID)
				handler.ServeHTTP(rr, req)
				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusNotFound))
			})
		})
		When("something bad happened", func() {
			fakeID, _ := faker.RandomInt(1000, 2000, 1)
			fakeIDUint := uint(fakeID[0])
			url := fmt.Sprintf("/%d", fakeIDUint)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			Expect(err).To(BeNil())

			It("should return status code 500", func() {
				ctx := req.Context()
				ctx = setContextDeviceGroup(ctx, &models.DeviceGroup{
					Model: models.Model{
						ID: fakeIDUint,
					},
				})
				ctx = dependencies.ContextWithServices(ctx, edgeAPIServices)
				req = req.WithContext(ctx)
				rr := httptest.NewRecorder()

				// setup mock for DeviceGroupsService
				mockDeviceGroupsService.EXPECT().DeleteDeviceGroupByID(fmt.Sprint(fakeIDUint)).Return(errors.NewInternalServerError())

				handler := http.HandlerFunc(DeleteDeviceGroupByID)
				handler.ServeHTTP(rr, req)
				// Check the status code is what we expect.
				Expect(rr.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
	Context("delete device from DeviceGroup", func() {
		account := faker.UUIDHyphenated()
		deviceGroupName := faker.Name()
		devices := []models.Device{
			{
				Name:    faker.Name(),
				UUID:    faker.UUIDHyphenated(),
				Account: account,
			},
			{
				Name:    faker.Name(),
				UUID:    faker.UUIDHyphenated(),
				Account: account,
			},
			{
				Name:    faker.Name(),
				UUID:    faker.UUIDHyphenated(),
				Account: account,
			},
		}
		deviceGroup := models.DeviceGroup{
			Name:    deviceGroupName,
			Account: account,
			Type:    models.DeviceGroupTypeDefault,
			Devices: devices,
		}

		It("should create device group with devices", func() {
			res := db.DB.Create(&deviceGroup)
			Expect(res.Error).To(BeNil())
			Expect(deviceGroup.ID).NotTo(Equal(0))
		})
		It("load device-group with devices", func() {
			res := db.DB.Preload("Devices").First(&deviceGroup, deviceGroup.ID)
			Expect(res.Error).To(BeNil())
			for _, device := range deviceGroup.Devices {
				// ensure all devices are defined
				Expect(device.ID).NotTo(Equal(0))
			}
		})
		When("device-group and devices are defined", func() {
			It("should delete the first device", func() {
				devicesToRemove := deviceGroup.Devices[:1]
				url := fmt.Sprintf("/%d/devices/%d", deviceGroup.ID, devicesToRemove[0].ID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				Expect(err).To(BeNil())

				ctx := req.Context()
				ctx = setContextDeviceGroup(ctx, &deviceGroup)
				ctx = setContextDeviceGroupDevice(ctx, &devicesToRemove[0])
				ctx = dependencies.ContextWithServices(ctx, edgeAPIServices)
				req = req.WithContext(ctx)
				rr := httptest.NewRecorder()

				mockDeviceGroupsService.EXPECT().DeleteDeviceGroupDevices(account, deviceGroup.ID, devicesToRemove).Return(&devicesToRemove, nil)
				handler := http.HandlerFunc(DeleteDeviceGroupOneDevice)
				handler.ServeHTTP(rr, req)
				Expect(rr.Code).To(Equal(http.StatusOK))
			})

			It("should delete the second and third devices", func() {

				devicesToRemove := deviceGroup.Devices[1:]
				devicesToRemoveJSON, err := json.Marshal(models.DeviceGroup{Devices: devicesToRemove})
				Expect(err).To(BeNil())

				url := fmt.Sprintf("/%d/devices", deviceGroup.ID)
				req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(devicesToRemoveJSON))
				Expect(err).To(BeNil())

				ctx := req.Context()
				ctx = setContextDeviceGroup(ctx, &deviceGroup)
				ctx = dependencies.ContextWithServices(ctx, edgeAPIServices)
				req = req.WithContext(ctx)
				rr := httptest.NewRecorder()

				mockDeviceGroupsService.EXPECT().DeleteDeviceGroupDevices(account, deviceGroup.ID, gomock.Any()).Return(&devicesToRemove, nil)
				handler := http.HandlerFunc(DeleteDeviceGroupManyDevices)
				handler.ServeHTTP(rr, req)

				Expect(rr.Code).To(Equal(http.StatusOK))
			})
		})
		When("sending invalid request body", func() {
			It("should return status code 400", func() {
				url := fmt.Sprintf("/%d/devices", deviceGroup.ID)
				req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer([]byte("{}")))
				Expect(err).To(BeNil())

				ctx := req.Context()
				ctx = setContextDeviceGroup(ctx, &deviceGroup)
				ctx = dependencies.ContextWithServices(ctx, edgeAPIServices)
				req = req.WithContext(ctx)
				rr := httptest.NewRecorder()

				var devicesToRemove []models.Device
				mockDeviceGroupsService.EXPECT().DeleteDeviceGroupDevices(account, deviceGroup.ID, devicesToRemove).Return(nil, new(services.DeviceGroupDevicesNotSupplied))
				handler := http.HandlerFunc(DeleteDeviceGroupManyDevices)
				handler.ServeHTTP(rr, req)
				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
