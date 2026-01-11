package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPPackingGateway(t *testing.T) {
	tests := []struct {
		name            string
		baseURL         string
		timeout         time.Duration
		expectedBaseURL string
	}{
		{
			name:            "with_valid_base_url",
			baseURL:         "http://packing-service:5051",
			timeout:         10 * time.Second,
			expectedBaseURL: "http://packing-service:5051",
		},
		{
			name:            "with_trailing_slash",
			baseURL:         "http://packing-service:5051/",
			timeout:         5 * time.Second,
			expectedBaseURL: "http://packing-service:5051",
		},
		{
			name:            "with_multiple_trailing_slashes",
			baseURL:         "http://packing-service:5051///",
			timeout:         5 * time.Second,
			expectedBaseURL: "http://packing-service:5051",
		},
		{
			name:            "empty_url_defaults_to_localhost",
			baseURL:         "",
			timeout:         5 * time.Second,
			expectedBaseURL: "http://localhost:5051",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gateway := NewHTTPPackingGateway(tt.baseURL, tt.timeout)

			assert.NotNil(t, gateway)
			assert.Equal(t, tt.expectedBaseURL, gateway.baseURL)
			assert.NotNil(t, gateway.client)
			assert.Equal(t, tt.timeout, gateway.client.Timeout)
		})
	}
}

func TestHTTPPackingGateway_Pack(t *testing.T) {
	tests := []struct {
		name           string
		request        PackRequest
		setupServer    func() *httptest.Server
		expectedError  string
		validateResult func(*testing.T, *PackResponse)
	}{
		{
			name: "successful_pack_with_placements",
			request: PackRequest{
				Units: "mm",
				Container: PackContainerIn{
					Length:    1000,
					Width:     800,
					Height:    600,
					MaxWeight: 5000,
				},
				Items: []PackItemIn{
					{
						ItemID:   "item-1",
						Label:    "Box A",
						Length:   200,
						Width:    150,
						Height:   100,
						Weight:   50,
						Quantity: 2,
					},
				},
				Options: PackOptionsIn{
					FixPoint:            true,
					CheckStable:         true,
					SupportSurfaceRatio: 0.75,
					BiggerFirst:         true,
					PutType:             1,
				},
			},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodPost, r.Method)
					assert.Equal(t, "/pack", r.URL.Path)
					assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

					response := PackResponse{
						Success: true,
						Data: &PackDataOut{
							Units: "mm",
							Placements: []PackPlacementOut{
								{
									ItemID:     "item-1",
									Label:      "Box A",
									PosX:       0,
									PosY:       0,
									PosZ:       0,
									Rotation:   0,
									StepNumber: 1,
								},
								{
									ItemID:     "item-1",
									Label:      "Box A",
									PosX:       200,
									PosY:       0,
									PosZ:       0,
									Rotation:   0,
									StepNumber: 2,
								},
							},
							Unfitted: []PackUnfittedOut{},
							Stats: PackStatsOut{
								ExpandedItems: 2,
								FittedCount:   2,
								UnfittedCount: 0,
								PackTimeMs:    15,
								TotalTimeMs:   20,
							},
						},
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				}))
			},
			validateResult: func(t *testing.T, resp *PackResponse) {
				assert.True(t, resp.Success)
				require.NotNil(t, resp.Data)
				assert.Equal(t, "mm", resp.Data.Units)
				assert.Len(t, resp.Data.Placements, 2)
				assert.Equal(t, "item-1", resp.Data.Placements[0].ItemID)
				assert.Equal(t, 1, resp.Data.Placements[0].StepNumber)
				assert.Equal(t, 2, resp.Data.Stats.FittedCount)
				assert.Equal(t, 0, resp.Data.Stats.UnfittedCount)
			},
		},
		{
			name: "partial_fit_with_unfitted_items",
			request: PackRequest{
				Units: "mm",
				Container: PackContainerIn{
					Length:    500,
					Width:     500,
					Height:    500,
					MaxWeight: 1000,
				},
				Items: []PackItemIn{
					{
						ItemID:   "large-box",
						Label:    "Large Box",
						Length:   600,
						Width:    600,
						Height:   600,
						Weight:   500,
						Quantity: 1,
					},
				},
				Options: PackOptionsIn{},
			},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := PackResponse{
						Success: true,
						Data: &PackDataOut{
							Units:      "mm",
							Placements: []PackPlacementOut{},
							Unfitted: []PackUnfittedOut{
								{
									ItemID: "large-box",
									Label:  "Large Box",
									Count:  1,
								},
							},
							Stats: PackStatsOut{
								ExpandedItems: 1,
								FittedCount:   0,
								UnfittedCount: 1,
								PackTimeMs:    5,
								TotalTimeMs:   8,
							},
						},
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				}))
			},
			validateResult: func(t *testing.T, resp *PackResponse) {
				assert.True(t, resp.Success)
				require.NotNil(t, resp.Data)
				assert.Len(t, resp.Data.Placements, 0)
				assert.Len(t, resp.Data.Unfitted, 1)
				assert.Equal(t, "large-box", resp.Data.Unfitted[0].ItemID)
				assert.Equal(t, 1, resp.Data.Unfitted[0].Count)
			},
		},
		{
			name:    "service_returns_non_2xx_status",
			request: PackRequest{Units: "mm"},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := PackResponse{
						Success: false,
						Error: &PackErrorOut{
							Code:    "INVALID_REQUEST",
							Message: "Container dimensions are required",
							Details: map[string]interface{}{
								"field": "container",
							},
						},
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(response)
				}))
			},
			expectedError: "packing service error: Container dimensions are required",
		},
		{
			name:    "service_returns_success_false",
			request: PackRequest{Units: "mm"},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := PackResponse{
						Success: false,
						Error: &PackErrorOut{
							Code:    "PACKING_FAILED",
							Message: "Unable to pack items",
						},
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				}))
			},
			expectedError: "packing service error: Unable to pack items",
		},
		{
			name:    "service_returns_success_true_but_missing_data",
			request: PackRequest{Units: "mm"},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := PackResponse{
						Success: true,
						Data:    nil,
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				}))
			},
			expectedError: "packing service error: missing data",
		},
		{
			name:    "service_returns_invalid_json",
			request: PackRequest{Units: "mm"},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("not valid json{"))
				}))
			},
			expectedError: "decode packing response (status 200)",
		},
		{
			name:    "service_returns_500_with_error_message",
			request: PackRequest{Units: "mm"},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := PackResponse{
						Success: false,
						Error: &PackErrorOut{
							Code:    "INTERNAL_ERROR",
							Message: "Database connection failed",
						},
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
				}))
			},
			expectedError: "packing service error: Database connection failed",
		},
		{
			name:    "service_returns_500_with_error_but_empty_message",
			request: PackRequest{Units: "mm"},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := PackResponse{
						Success: false,
						Error: &PackErrorOut{
							Code:    "INTERNAL_ERROR",
							Message: "", // Empty message
						},
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
				}))
			},
			expectedError: "packing service error: packing service returned 500",
		},
		{
			name:    "service_returns_500_without_error_object",
			request: PackRequest{Units: "mm"},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := PackResponse{
						Success: false,
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
				}))
			},
			expectedError: "packing service error: packing service returned 500",
		},
		{
			name:    "service_returns_success_false_with_error_but_empty_message",
			request: PackRequest{Units: "mm"},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := PackResponse{
						Success: false,
						Error: &PackErrorOut{
							Code:    "PACKING_FAILED",
							Message: "", // Empty message
						},
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				}))
			},
			expectedError: "packing service error: packing service returned success=false",
		},
		{
			name:    "service_returns_success_false_without_error_object",
			request: PackRequest{Units: "mm"},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := PackResponse{
						Success: false,
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				}))
			},
			expectedError: "packing service error: packing service returned success=false",
		},
		{
			name:    "network_error_service_unreachable",
			request: PackRequest{Units: "mm"},
			setupServer: func() *httptest.Server {
				// Create server but immediately close it
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
				server.Close()
				return server
			},
			expectedError: "call packing service",
		},
		{
			name:    "context_timeout",
			request: PackRequest{Units: "mm"},
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Simulate slow server
					time.Sleep(200 * time.Millisecond)
					w.WriteHeader(http.StatusOK)
				}))
			},
			expectedError: "call packing service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			gateway := NewHTTPPackingGateway(server.URL, 5*time.Second)

			ctx := context.Background()
			if tt.name == "context_timeout" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, 50*time.Millisecond)
				defer cancel()
			}

			result, err := gateway.Pack(ctx, tt.request)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}
		})
	}
}

func TestHTTPPackingGateway_Pack_RequestValidation(t *testing.T) {
	// Test that request is properly marshaled and sent
	t.Run("request_body_marshaled_correctly", func(t *testing.T) {
		var receivedRequest PackRequest

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Decode the request to verify it was sent correctly
			err := json.NewDecoder(r.Body).Decode(&receivedRequest)
			require.NoError(t, err)

			response := PackResponse{
				Success: true,
				Data: &PackDataOut{
					Units:      "mm",
					Placements: []PackPlacementOut{},
					Unfitted:   []PackUnfittedOut{},
					Stats:      PackStatsOut{},
				},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		gateway := NewHTTPPackingGateway(server.URL, 5*time.Second)

		request := PackRequest{
			Units: "mm",
			Container: PackContainerIn{
				Length:    1000,
				Width:     800,
				Height:    600,
				MaxWeight: 5000,
			},
			Items: []PackItemIn{
				{
					ItemID:   "test-item",
					Label:    "Test Item",
					Length:   100,
					Width:    100,
					Height:   100,
					Weight:   10,
					Quantity: 5,
				},
			},
			Options: PackOptionsIn{
				FixPoint:            true,
				CheckStable:         false,
				SupportSurfaceRatio: 0.8,
				BiggerFirst:         true,
				PutType:             2,
			},
		}

		_, err := gateway.Pack(context.Background(), request)
		require.NoError(t, err)

		// Verify the request was received correctly
		assert.Equal(t, "mm", receivedRequest.Units)
		assert.Equal(t, 1000.0, receivedRequest.Container.Length)
		assert.Equal(t, 800.0, receivedRequest.Container.Width)
		assert.Equal(t, 600.0, receivedRequest.Container.Height)
		assert.Equal(t, 5000.0, receivedRequest.Container.MaxWeight)
		require.Len(t, receivedRequest.Items, 1)
		assert.Equal(t, "test-item", receivedRequest.Items[0].ItemID)
		assert.Equal(t, "Test Item", receivedRequest.Items[0].Label)
		assert.Equal(t, 5, receivedRequest.Items[0].Quantity)
		assert.True(t, receivedRequest.Options.FixPoint)
		assert.False(t, receivedRequest.Options.CheckStable)
		assert.Equal(t, 0.8, receivedRequest.Options.SupportSurfaceRatio)
		assert.Equal(t, 2, receivedRequest.Options.PutType)
	})
}

func TestHTTPPackingGateway_Pack_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		setupServer   func() *httptest.Server
		expectedError string
	}{
		{
			name: "empty_response_body",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					// Send empty body
				}))
			},
			expectedError: "decode packing response",
		},
		{
			name: "response_with_null_error",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					response := PackResponse{
						Success: true,
						Data: &PackDataOut{
							Units:      "mm",
							Placements: []PackPlacementOut{},
							Unfitted:   []PackUnfittedOut{},
							Stats:      PackStatsOut{},
						},
						Error: nil,
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(response)
				}))
			},
			expectedError: "",
		},
		{
			name: "response_body_read_error",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("Content-Length", "1000") // Claim large body
					w.WriteHeader(http.StatusOK)
					// Close connection immediately without sending body
					// This will cause io.ReadAll to potentially fail
					hj, ok := w.(http.Hijacker)
					if ok {
						conn, _, _ := hj.Hijack()
						conn.Close()
					}
				}))
			},
			expectedError: "read packing response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			gateway := NewHTTPPackingGateway(server.URL, 5*time.Second)

			request := PackRequest{Units: "mm"}
			result, err := gateway.Pack(context.Background(), request)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
			}
		})
	}
}
