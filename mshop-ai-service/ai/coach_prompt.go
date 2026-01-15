package ai

import (
	"encoding/json"
	"strings"
)

var allowedCoachIntents = map[string]bool{
	"RECOVERY_HINT_GET": true,
	"WANTS_INFO":        true,
	"UNKNOWN":           true,
}

type coachOut struct {
	Intent string                 `json:"intent"`
	Params map[string]interface{} `json:"params,omitempty"`
}

func BuildCoachPrompt(userMessage string) string {
	system := `
You are an AI helpdesk assistant for a mobile app.
Return ONLY valid JSON and nothing else. No markdown, no explanations.

Allowed intents:
- "RECOVERY_HINT_GET"  (user asks where recovery code is / wants reminder)
- "WANTS_INFO"         (message unclear: ask a short follow-up / guide)
- "UNKNOWN"            (not understood)

Rules:
1) If user asks where they saved the recovery code -> {"intent":"RECOVERY_HINT_GET"}
2) If user has login trouble but doesn't explicitly ask for recovery code ->
   {"intent":"WANTS_INFO","params":{"message":"Trebaš li podsjetnik gdje si spremila recovery kod? Ako da, napiši: 'Gdje mi je recovery kod?'" }}
3) If not understood -> {"intent":"UNKNOWN"}

When using WANTS_INFO, respond in clear, natural Croatian.
Do NOT use literal translations.
Use short, polite, user-friendly sentences suitable for a mobile app.
Examples:
- "Nisam siguran što točno tražiš. Možeš li malo pojasniti?"
- "Ako tražiš recovery kod, napiši: 'Gdje mi je recovery kod?'"

`
	return system + "\nUser message: " + `"` + userMessage + `"`
}

func ExtractAndValidateCoachJSON(raw string) string {
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start == -1 || end == -1 || end < start {
		return `{"intent":"UNKNOWN"}`
	}

	jsonPart := raw[start : end+1]

	var out coachOut
	if err := json.Unmarshal([]byte(jsonPart), &out); err != nil {
		return `{"intent":"UNKNOWN"}`
	}

	if !allowedCoachIntents[out.Intent] {
		return `{"intent":"UNKNOWN"}`
	}

	b, err := json.Marshal(out)
	if err != nil {
		return `{"intent":"UNKNOWN"}`
	}
	return string(b)
}
