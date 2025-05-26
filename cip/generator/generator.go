package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// GenerateCIP generates the complete CIP implementation from a definition
func GenerateCIP(def CIPDefinition) error {
	// Create the package directory if it doesn't exist
	dir := filepath.Join("cip", def.Package)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Update constants and KeyOpMap
	if err := updateConstants(def); err != nil {
		return fmt.Errorf("failed to update constants: %v", err)
	}

	if err := updateKeyOpMap(def); err != nil {
		return fmt.Errorf("failed to update KeyOpMap: %v", err)
	}

	// Generate the main implementation file
	if err := generateImplementationFile(def); err != nil {
		return err
	}

	// Generate the test file
	if err := generateTestFile(def); err != nil {
		return err
	}

	return nil
}

// pascalCase returns a PascalCase version of the input string
func pascalCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// updateConstants updates the constants.go file with new event kinds and operations
func updateConstants(def CIPDefinition) error {
	constantsFile := filepath.Join("cip", "constants.go")
	content, err := os.ReadFile(constantsFile)
	if err != nil {
		return fmt.Errorf("failed to read constants file: %v", err)
	}

	// Find all const blocks
	constBlocks := []struct {
		startMarker string
		endMarker   string
		comment     string
		format      string
	}{
		{
			startMarker: "// Event kinds",
			endMarker:   ")\n\n// Event operations",
			comment:     def.CIPName + " event kinds",
			format:      "\tKind%s = %d\n",
		},
		{
			startMarker: "// Event operations",
			endMarker:   ")\n\n// Default cip operations",
			comment:     def.CIPName + " operation types",
			format:      "\tOp%s = \"%s\" // %d\n",
		},
		{
			startMarker: "// Default cip operations",
			endMarker:   ")\n",
			comment:     def.CIPName + " operations string",
			format:      "\t%s = \"%s\"\n",
		},
	}

	// Prepare new content
	newContent := string(content)
	for _, block := range constBlocks {
		// Find the start of the block
		startIndex := strings.Index(newContent, block.startMarker)
		if startIndex == -1 {
			return fmt.Errorf("could not find %s block", block.startMarker)
		}

		// Find the end of the block
		endIndex := strings.Index(newContent[startIndex:], block.endMarker)
		if endIndex == -1 {
			return fmt.Errorf("could not find end of %s block", block.startMarker)
		}
		endIndex += startIndex

		// Prepare new constants
		var newConstants strings.Builder
		newConstants.WriteString("\n\t// " + block.comment + "\n")

		switch block.startMarker {
		case "// Event kinds":
			for _, event := range def.Events {
				name := strings.TrimSuffix(event.EventName, "Event")
				pascalCIP := pascalCase(def.CIPName)
				pascalName := pascalCase(name)
				newConstants.WriteString(fmt.Sprintf(block.format, pascalCIP+pascalName, event.Kind))
			}
		case "// Event operations":
			for _, event := range def.Events {
				name := strings.TrimSuffix(event.EventName, "Event")
				pascalName := pascalCase(name)
				newConstants.WriteString(fmt.Sprintf(block.format, pascalName, event.Operation, event.Kind))
			}
		case "// Default cip operations":
			// Build operations string
			var ops []string
			for _, event := range def.Events {
				ops = append(ops, fmt.Sprintf("%s=%d", event.Operation, event.Kind))
			}
			pascalCIP := pascalCase(def.CIPName)
			newConstants.WriteString(fmt.Sprintf(block.format, pascalCIP+"SubspaceOps", strings.Join(ops, ",")))
		}

		// Insert new constants before the end of the block
		newContent = newContent[:endIndex] + newConstants.String() + newContent[endIndex:]
	}

	// Write back to file
	if err := os.WriteFile(constantsFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write constants file: %v", err)
	}

	return nil
}

// updateKeyOpMap updates the keys.go file with new operation mappings
func updateKeyOpMap(def CIPDefinition) error {
	keysFile := filepath.Join("cip", "keys.go")
	content, err := os.ReadFile(keysFile)
	if err != nil {
		return fmt.Errorf("failed to read keys file: %v", err)
	}

	// Find the start of the KeyOpMap
	mapStart := strings.Index(string(content), "var KeyOpMap = map[int]string{")
	if mapStart == -1 {
		return fmt.Errorf("could not find KeyOpMap")
	}

	// Find the closing } of the map
	mapEnd := strings.Index(string(content[mapStart:]), "}")
	if mapEnd == -1 {
		return fmt.Errorf("could not find end of KeyOpMap")
	}
	mapEnd += mapStart

	// Prepare new operation mappings
	var newOps strings.Builder
	newOps.WriteString("\n\t// " + def.CIPName + " operations\n")
	for _, event := range def.Events {
		name := strings.TrimSuffix(event.EventName, "Event")
		pascalCIP := pascalCase(def.CIPName)
		pascalName := pascalCase(name)
		newOps.WriteString(fmt.Sprintf("\tKind%s%s: Op%s,\n", pascalCIP, pascalName, pascalName))
	}

	// Insert new operations before the closing }
	newContent := string(content[:mapEnd]) + newOps.String() + string(content[mapEnd:])

	// Write back to file
	if err := os.WriteFile(keysFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write keys file: %v", err)
	}

	return nil
}

const eventTemplate = `package {{.Package}}

import (
	"fmt"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
)

{{- range .Events}}
{{$event := .}}
// {{.EventName}} represents a {{.Operation}} operation in {{$.CIPName}} subspace
type {{.EventName}} struct {
	*nostr.SubspaceOpEvent
{{- range .Fields}}
	{{.FieldName}} {{if .Multiple}}[]{{end}}{{.Type}}
{{- end}}
}

{{- range .Fields}}

// Set{{.FieldName}} sets the {{.FieldName}} for the operation
func (e *{{$event.EventName}}) Set{{.FieldName}}(value {{if .Multiple}}[]{{end}}{{.Type}}) {
{{- if .Multiple}}
	e.{{.FieldName}} = value
	tag := nostr.Tag{"{{.Tag}}"}
	for _, v := range value {
		tag = append(tag, v)
	}
	e.Tags = append(e.Tags, tag)
{{- else}}
	e.{{.FieldName}} = value
	e.Tags = append(e.Tags, nostr.Tag{"{{.Tag}}", value})
{{- end}}
}

{{- end}}

// New{{.EventName}} creates a new {{.Operation}} event
func New{{.EventName}}(subspaceID string) (*{{.EventName}}, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.Kind{{pascalCase $.CIPName}}{{pascalCase (trimSuffix .EventName "Event")}})
	if err != nil {
		return nil, err
	}
	return &{{.EventName}}{
		SubspaceOpEvent: baseEvent,
	}, nil
}

{{- end}}

// Parse{{pascalCase .CIPName}}Event parses a Nostr event into a {{.CIPName}} event
func Parse{{pascalCase .CIPName}}Event(evt nostr.Event) (nostr.SubspaceOpEventPtr, error) {
	subspaceID := ""
	parents := []string{}
	var authTag cip.AuthTag
	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "sid":
			subspaceID = tag[1]
		case "auth":
			auth, err := cip.ParseAuthTag(tag[1])
			if err != nil {
				return nil, fmt.Errorf("failed to parse auth tag: %v", err)
			}
			authTag = auth
		case "parent":
			parents = append(parents, tag[1:]...)
		}
	}
	operation, exists := cip.GetOpFromKind(evt.Kind)
	if !exists {
		return nil, fmt.Errorf("unknown kind value: %d", evt.Kind)
	}
	switch operation {
{{- range .Events}}
	case cip.Op{{pascalCase (trimSuffix .EventName "Event")}}:
		return parse{{.EventName}}(evt, subspaceID, operation, authTag, parents)
{{- end}}
	default:
		return nil, fmt.Errorf("unknown operation type: %s", operation)
	}
}

{{- range .Events}}

func parse{{.EventName}}(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*{{.EventName}}, error) {
	event := &{{.EventName}}{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
{{- range .Fields}}
		case "{{.Tag}}":
{{- if .Multiple}}
			event.{{.FieldName}} = tag[1:]
{{- else}}
			event.{{.FieldName}} = tag[1]
{{- end}}
{{- end}}
		}
	}
	return event, nil
}

{{- end}}
`

func generateImplementationFile(def CIPDefinition) error {
	tmpl := eventTemplate

	t, err := template.New("implementation").Funcs(template.FuncMap{
		"trimSuffix": strings.TrimSuffix,
		"pascalCase": pascalCase,
	}).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	f, err := os.Create(filepath.Join("cip", def.Package, strings.ToLower(def.CIPName)+".go"))
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()

	if err := t.Execute(f, def); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}

func generateTestFile(def CIPDefinition) error {
	tmpl := `package {{.Package}}

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/nbd-wtf/go-nostr/cip"
)

{{- range .Events}}
func Test{{.EventName}}(t *testing.T) {
	// Test creating a new event
	event, err := New{{.EventName}}("test-subspace")
	assert.NoError(t, err)
	assert.Equal(t, cip.Kind{{pascalCase $.CIPName}}{{pascalCase (trimSuffix .EventName "Event")}}, event.Kind)

	{{- range .Fields}}
	{{- if .Multiple}}
	// Test setting {{.FieldName}}
	values := []{{.Type}}{"value1", "value2"}
	event.Set{{.FieldName}}(values)
	assert.Equal(t, values, event.{{.FieldName}})
	{{- else}}
	// Test setting {{.FieldName}}
	event.Set{{.FieldName}}("test-value")
	assert.Equal(t, "test-value", event.{{.FieldName}})
	{{- end}}
	{{- end}}

	// Test parsing
	parsedEvent, err := Parse{{pascalCase $.CIPName}}Event(event.Event)
	assert.NoError(t, err)
	assert.IsType(t, &{{.EventName}}{}, parsedEvent)
}
{{- end}}
`

	t, err := template.New("test").Funcs(template.FuncMap{
		"trimSuffix": strings.TrimSuffix,
		"pascalCase": pascalCase,
	}).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	f, err := os.Create(filepath.Join("cip", def.Package, strings.ToLower(def.CIPName)+"_test.go"))
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()

	if err := t.Execute(f, def); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return nil
}
