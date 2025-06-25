package proxy

import (
	"encoding/json"
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
			content, _ := msgMap["content"].(string)
			
			if role == "system" {
				systemContents = append(systemContents, content)
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