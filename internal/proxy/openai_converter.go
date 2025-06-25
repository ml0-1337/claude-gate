package proxy

import (
	"encoding/json"
	"log/slog"
	"strings"
	"time"
)

// ConvertOpenAIToAnthropic converts OpenAI chat/completions format to Anthropic messages format
func ConvertOpenAIToAnthropic(body []byte) ([]byte, error) {
	var openAIRequest map[string]interface{}
	if err := json.Unmarshal(body, &openAIRequest); err != nil {
		return nil, err
	}
	
	// Create Anthropic format request
	anthropicRequest := make(map[string]interface{})
	
	// Convert model name (remove "anthropic/" prefix if present)
	if model, ok := openAIRequest["model"].(string); ok {
		anthropicRequest["model"] = strings.TrimPrefix(model, "anthropic/")
	}
	
	// Extract system messages and convert messages array
	var systemContents []string
	var anthropicMessages []interface{}
	
	if messages, ok := openAIRequest["messages"].([]interface{}); ok {
		for _, msg := range messages {
			msgMap, ok := msg.(map[string]interface{})
			if !ok {
				continue
			}
			
			role, _ := msgMap["role"].(string)
			
			// Handle content - can be string or array
			var content interface{}
			if msgContent, ok := msgMap["content"]; ok {
				switch v := msgContent.(type) {
				case string:
					content = v
				case []interface{}:
					// Handle structured content array
					content = v
				default:
					continue
				}
			}
			
			if role == "system" {
				// Extract text from system message
				switch v := content.(type) {
				case string:
					systemContents = append(systemContents, v)
				case []interface{}:
					// Handle array of content items
					for _, item := range v {
						if itemMap, ok := item.(map[string]interface{}); ok {
							if itemMap["type"] == "text" {
								if text, ok := itemMap["text"].(string); ok {
									systemContents = append(systemContents, text)
								}
							}
						}
					}
				}
			} else {
				// Convert to Anthropic message format
				anthropicMsg := map[string]interface{}{
					"role":    role,
					"content": content,
				}
				anthropicMessages = append(anthropicMessages, anthropicMsg)
			}
		}
	}
	
	// Set messages
	anthropicRequest["messages"] = anthropicMessages
	
	// Build system field with Claude Code prompt first
	systemArray := []interface{}{
		map[string]interface{}{"type": "text", "text": ClaudeCodePrompt},
	}
	
	// Add extracted system messages
	for _, systemContent := range systemContents {
		systemArray = append(systemArray, map[string]interface{}{
			"type": "text",
			"text": systemContent,
		})
	}
	
	anthropicRequest["system"] = systemArray
	
	// Copy other fields
	for key, value := range openAIRequest {
		if key != "model" && key != "messages" {
			anthropicRequest[key] = value
		}
	}
	
	return json.Marshal(anthropicRequest)
}

// ConvertAnthropicToOpenAI converts Anthropic response format to OpenAI chat/completions format
func ConvertAnthropicToOpenAI(body []byte) ([]byte, error) {
	var anthropicResponse map[string]interface{}
	if err := json.Unmarshal(body, &anthropicResponse); err != nil {
		return nil, err
	}
	
	// Check if this is an error response
	if errorObj, hasError := anthropicResponse["error"]; hasError {
		// Convert Anthropic error to OpenAI error format
		return convertAnthropicErrorToOpenAI(errorObj)
	}
	
	// Create OpenAI format response
	openAIResponse := make(map[string]interface{})
	
	// Copy basic fields
	if id, ok := anthropicResponse["id"].(string); ok {
		openAIResponse["id"] = id
	}
	openAIResponse["object"] = "chat.completion"
	if model, ok := anthropicResponse["model"].(string); ok {
		openAIResponse["model"] = model
	}
	openAIResponse["created"] = int(time.Now().Unix())
	
	// Convert content to OpenAI format
	var messageContent string
	if content, ok := anthropicResponse["content"].([]interface{}); ok {
		for _, item := range content {
			if contentMap, ok := item.(map[string]interface{}); ok {
				if contentMap["type"] == "text" {
					if text, ok := contentMap["text"].(string); ok {
						messageContent += text
					}
				}
			}
		}
	}
	
	// Build choices array
	finishReason := "stop"
	if stopReason, ok := anthropicResponse["stop_reason"].(string); ok {
		switch stopReason {
		case "end_turn":
			finishReason = "stop"
		case "max_tokens":
			finishReason = "length"
		case "stop_sequence":
			finishReason = "stop"
		default:
			finishReason = stopReason
		}
	}
	
	choices := []interface{}{
		map[string]interface{}{
			"index": 0,
			"message": map[string]interface{}{
				"role":    "assistant",
				"content": messageContent,
			},
			"finish_reason": finishReason,
		},
	}
	openAIResponse["choices"] = choices
	
	// Convert usage
	if anthropicUsage, ok := anthropicResponse["usage"].(map[string]interface{}); ok {
		inputTokens := 0
		outputTokens := 0
		
		if val, ok := anthropicUsage["input_tokens"].(float64); ok {
			inputTokens = int(val)
		}
		if val, ok := anthropicUsage["output_tokens"].(float64); ok {
			outputTokens = int(val)
		}
		
		openAIResponse["usage"] = map[string]interface{}{
			"prompt_tokens":     inputTokens,
			"completion_tokens": outputTokens,
			"total_tokens":      inputTokens + outputTokens,
		}
	}
	
	return json.Marshal(openAIResponse)
}

// convertAnthropicErrorToOpenAI converts Anthropic error format to OpenAI error format
func convertAnthropicErrorToOpenAI(errorObj interface{}) ([]byte, error) {
	openAIError := map[string]interface{}{
		"error": map[string]interface{}{
			"message": "",
			"type":    "invalid_request_error",
			"param":   nil,
			"code":    nil,
		},
	}
	
	// Extract error details from Anthropic error
	if errorMap, ok := errorObj.(map[string]interface{}); ok {
		if message, ok := errorMap["message"].(string); ok {
			openAIError["error"].(map[string]interface{})["message"] = message
		}
		
		// Map Anthropic error types to OpenAI error types
		if errorType, ok := errorMap["type"].(string); ok {
			switch errorType {
			case "invalid_request_error":
				openAIError["error"].(map[string]interface{})["type"] = "invalid_request_error"
			case "authentication_error":
				openAIError["error"].(map[string]interface{})["type"] = "authentication_error"
			case "permission_error":
				openAIError["error"].(map[string]interface{})["type"] = "permission_denied"
			case "not_found_error":
				openAIError["error"].(map[string]interface{})["type"] = "not_found_error"
			case "rate_limit_error":
				openAIError["error"].(map[string]interface{})["type"] = "rate_limit_error"
			case "api_error":
				openAIError["error"].(map[string]interface{})["type"] = "server_error"
			default:
				openAIError["error"].(map[string]interface{})["type"] = errorType
			}
		}
	}
	
	return json.Marshal(openAIError)
}

// toolState tracks tool use information across SSE events
// Key is Anthropic content block index, value contains tool info and OpenAI tool index
var toolState = make(map[int]map[string]interface{})

// toolCallIndex tracks the next available OpenAI tool call index
var toolCallIndex = 0

// ResetSSEConverterState resets the converter state (useful for testing)
func ResetSSEConverterState() {
	toolState = make(map[int]map[string]interface{})
	toolCallIndex = 0
}

// ConvertAnthropicSSEToOpenAI converts a single Anthropic SSE event to OpenAI format
func ConvertAnthropicSSEToOpenAI(event, data string, messageID string, model string, created int64) (string, error) {
	return ConvertAnthropicSSEToOpenAIWithLogger(event, data, messageID, model, created, nil)
}

// ConvertAnthropicSSEToOpenAIWithLogger converts a single Anthropic SSE event to OpenAI format with optional logging
func ConvertAnthropicSSEToOpenAIWithLogger(event, data string, messageID string, model string, created int64, logger *slog.Logger) (string, error) {
	// Parse the data as JSON
	var eventData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &eventData); err != nil {
		return "", err
	}
	
	eventType, _ := eventData["type"].(string)
	
	switch eventType {
	case "message_start":
		// Reset tool call index for new message
		toolCallIndex = 0
		
		// Convert message_start to initial OpenAI chunk
		chunk := map[string]interface{}{
			"id":      messageID,
			"object":  "chat.completion.chunk",
			"created": created,
			"model":   model,
			"choices": []interface{}{
				map[string]interface{}{
					"index": 0,
					"delta": map[string]interface{}{
						"role": "assistant",
					},
					"finish_reason": nil,
				},
			},
		}
		chunkJSON, _ := json.Marshal(chunk)
		return "data: " + string(chunkJSON) + "\n\n", nil
		
	case "content_block_start":
		// Check if this is a tool use block
		if contentBlock, ok := eventData["content_block"].(map[string]interface{}); ok {
			if contentBlock["type"] == "tool_use" {
				// Extract tool information
				index := int(eventData["index"].(float64))
				toolID, _ := contentBlock["id"].(string)
				toolName, _ := contentBlock["name"].(string)
				
				// Store tool info for this index with OpenAI tool index
				currentToolIndex := toolCallIndex
				toolState[index] = map[string]interface{}{
					"id":        toolID,
					"name":      toolName,
					"toolIndex": currentToolIndex,
				}
				toolCallIndex++
				
				// Send initial tool call chunk
				chunk := map[string]interface{}{
					"id":      messageID,
					"object":  "chat.completion.chunk",
					"created": created,
					"model":   model,
					"choices": []interface{}{
						map[string]interface{}{
							"index": 0,
							"delta": map[string]interface{}{
								"tool_calls": []interface{}{
									map[string]interface{}{
										"index": currentToolIndex,
										"id":    toolID,
										"type":  "function",
										"function": map[string]interface{}{
											"name":      toolName,
											"arguments": "",
										},
									},
								},
							},
							"finish_reason": nil,
						},
					},
				}
				chunkJSON, _ := json.Marshal(chunk)
				return "data: " + string(chunkJSON) + "\n\n", nil
			}
		}
		// For non-tool blocks, return empty to skip
		return "", nil
		
	case "content_block_stop":
		// Clear tool state for the completed block if it was a tool block
		if eventData["index"] != nil {
			index := int(eventData["index"].(float64))
			if _, exists := toolState[index]; exists {
				delete(toolState, index)
				if logger != nil {
					logger.Debug("completed tool use block", "index", index)
				}
			}
		}
		// No output for content_block_stop
		return "", nil
		
	case "content_block_delta":
		// Convert content delta to OpenAI chunk
		if delta, ok := eventData["delta"].(map[string]interface{}); ok {
			if delta["type"] == "text_delta" {
				if text, ok := delta["text"].(string); ok {
					chunk := map[string]interface{}{
						"id":      messageID,
						"object":  "chat.completion.chunk",
						"created": created,
						"model":   model,
						"choices": []interface{}{
							map[string]interface{}{
								"index": 0,
								"delta": map[string]interface{}{
									"content": text,
								},
								"finish_reason": nil,
							},
						},
					}
					chunkJSON, _ := json.Marshal(chunk)
					return "data: " + string(chunkJSON) + "\n\n", nil
				}
			} else if delta["type"] == "input_json_delta" {
				// Handle tool use deltas
				if partialJSON, ok := delta["partial_json"].(string); ok {
					// Get the content block index
					blockIndex := int(eventData["index"].(float64))
					
					// Look up the tool info for this block
					if toolInfo, exists := toolState[blockIndex]; exists {
						toolIndex := toolInfo["toolIndex"].(int)
						toolID := toolInfo["id"].(string)
						
						// Create tool delta chunk with tool ID
						chunk := map[string]interface{}{
							"id":      messageID,
							"object":  "chat.completion.chunk",
							"created": created,
							"model":   model,
							"choices": []interface{}{
								map[string]interface{}{
									"index": 0,
									"delta": map[string]interface{}{
										"tool_calls": []interface{}{
											map[string]interface{}{
												"index": toolIndex,
												"id":    toolID,
												"function": map[string]interface{}{
													"arguments": partialJSON,
												},
											},
										},
									},
									"finish_reason": nil,
								},
							},
						}
						chunkJSON, _ := json.Marshal(chunk)
						return "data: " + string(chunkJSON) + "\n\n", nil
					}
				}
			}
		}
		
	case "message_stop":
		// Send final chunk with finish_reason
		// Note: [DONE] marker should be sent separately by the stream handler
		chunk := map[string]interface{}{
			"id":      messageID,
			"object":  "chat.completion.chunk",
			"created": created,
			"model":   model,
			"choices": []interface{}{
				map[string]interface{}{
					"index":         0,
					"delta":         map[string]interface{}{},
					"finish_reason": "stop",
				},
			},
		}
		chunkJSON, _ := json.Marshal(chunk)
		return "data: " + string(chunkJSON) + "\n\n", nil
		
	case "message_delta":
		// Handle stop reasons from message_delta
		if delta, ok := eventData["delta"].(map[string]interface{}); ok {
			if stopReason, ok := delta["stop_reason"].(string); ok {
				finishReason := "stop"
				if stopReason == "max_tokens" {
					finishReason = "length"
				}
				
				chunk := map[string]interface{}{
					"id":      messageID,
					"object":  "chat.completion.chunk",
					"created": created,
					"model":   model,
					"choices": []interface{}{
						map[string]interface{}{
							"index":         0,
							"delta":         map[string]interface{}{},
							"finish_reason": finishReason,
						},
					},
				}
				chunkJSON, _ := json.Marshal(chunk)
				return "data: " + string(chunkJSON) + "\n\n", nil
			}
		}
		
	default:
		// Log unhandled event types for debugging
		if logger != nil {
			logger.Debug("unhandled SSE event type",
				"event", event,
				"type", eventType,
				"data", data,
			)
		}
	}
	
	// Skip other event types
	return "", nil
}