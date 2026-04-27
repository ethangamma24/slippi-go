package realtime

import nativeslippi "github.com/ethangamma24/slippi-go/pkg/slippi/types"

// LatestFrame returns the highest numbered finalized frame.
func LatestFrame(frames map[int]nativeslippi.Frame) (nativeslippi.Frame, bool) {
	latest := -1 << 30
	var out nativeslippi.Frame
	for idx, frame := range frames {
		if idx > latest {
			latest = idx
			out = frame
		}
	}
	if latest == -1<<30 {
		return nativeslippi.Frame{}, false
	}
	return out, true
}
