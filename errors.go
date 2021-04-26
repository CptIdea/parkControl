package main

var (
	ERRFailAuth     = "{\"data\":[{\"error\":{\"id\":1,\"text\":\"Auth failed.\"}}],\"result\":\"error\"}"
	ERRParkFull     = "{\"data\":[{\"error\":{\"id\":2,\"text\":\"Parking failed.\"}}],\"result\":\"error\"}"
	ERRSlotNotFound = "{\"data\":[{\"error\":{\"id\":3,\"text\":\"Slot not found.\"}}],\"result\":\"error\"}"
	ERRMapNotFound  = "{\"data\":[{\"error\":{\"id\":4,\"text\":\"Map not found.\"}}],\"result\":\"error\"}"
	ERRUserCanT     = "{\"data\":[{\"error\":{\"id\":5,\"text\":\"The action is impossible.\"}}],\"result\":\"error\"}"
	ERRSlotBlocked  = "{\"data\":[{\"error\":{\"id\":6,\"text\":\"Slot blocked.\"}}],\"result\":\"error\"}"
	SendOK          = "{\"result\":\"ok\",\"data\":[]}"
)
