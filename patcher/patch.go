package patcher

import "strings"

// A Patch is the manifest's representation of a patch.
type Patch struct {
	Name    string
	Patcher Patcher
}

// NewPatchFromMap creates a Patch from a map.
func NewPatchFromMap(mp map[string]interface{}) *Patch {
	patch := new(Patch)

	source, _ := mp["source"].(string)
	destination, _ := mp["destination"].(string)

	// Name
	tokens := strings.Split(source, "/")
	patch.Name = tokens[len(tokens)-1]

	switch mp["type"] {
	case "file":
		fallthrough
	default:
		patch.Patcher = NewFilePatcher(source, destination)
	}

	return patch
}

// PatchResult is the result of a patch operation.
type PatchResult struct {
	Error   error
	Patcher Patcher
}

// PatchResults is the result of a set of patches
type PatchResults []*PatchResult

// GetSuccess gets the successful patches.
func (pr PatchResults) GetSuccess() PatchResults {
	var successes PatchResults

	for _, r := range pr {
		if r.Error == nil {
			successes = append(successes, r)
		}
	}

	return successes
}

// GetFailed gets the failed patches.
func (pr PatchResults) GetFailed() PatchResults {
	var failed PatchResults

	for _, r := range pr {
		if r.Error != nil {
			failed = append(failed, r)
		}
	}

	return failed
}
