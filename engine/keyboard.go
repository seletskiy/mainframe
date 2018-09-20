package engine

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

func MapKeyToSymbol(key glfw.Key, mods glfw.ModifierKey) string {
	symbol := InputKeys[key]
	if symbol == "" {
		symbol = InputKeyUnknown
	}

	ignore := map[glfw.Key]bool{
		glfw.KeyLeftControl:  true,
		glfw.KeyLeftShift:    true,
		glfw.KeyLeftAlt:      true,
		glfw.KeyLeftSuper:    true,
		glfw.KeyRightControl: true,
		glfw.KeyRightShift:   true,
		glfw.KeyRightAlt:     true,
		glfw.KeyRightSuper:   true,
	}

	if ignore[key] {
		return symbol
	}

	symbol = InputMods[mods&glfw.ModControl] + symbol
	symbol = InputMods[mods&glfw.ModShift] + symbol
	symbol = InputMods[mods&glfw.ModAlt] + symbol
	symbol = InputMods[mods&glfw.ModSuper] + symbol

	return symbol
}

var InputKeyUnknown = `Unknown`

var InputMods = map[glfw.ModifierKey]string{
	glfw.ModControl: `Ctrl+`,
	glfw.ModShift:   `Shift+`,
	glfw.ModAlt:     `Alt+`,
	glfw.ModSuper:   `Super+`,
}

var InputKeys = map[glfw.Key]string{
	glfw.KeySpace:        `Space`,
	glfw.KeyApostrophe:   `'`,
	glfw.KeyComma:        `,`,
	glfw.KeyMinus:        `-`,
	glfw.KeyPeriod:       `.`,
	glfw.KeySlash:        `/`,
	glfw.Key0:            `0`,
	glfw.Key1:            `1`,
	glfw.Key2:            `2`,
	glfw.Key3:            `3`,
	glfw.Key4:            `4`,
	glfw.Key5:            `5`,
	glfw.Key6:            `6`,
	glfw.Key7:            `7`,
	glfw.Key8:            `8`,
	glfw.Key9:            `9`,
	glfw.KeySemicolon:    `;`,
	glfw.KeyEqual:        `=`,
	glfw.KeyA:            `A`,
	glfw.KeyB:            `B`,
	glfw.KeyC:            `C`,
	glfw.KeyD:            `D`,
	glfw.KeyE:            `E`,
	glfw.KeyF:            `F`,
	glfw.KeyG:            `G`,
	glfw.KeyH:            `H`,
	glfw.KeyI:            `I`,
	glfw.KeyJ:            `J`,
	glfw.KeyK:            `K`,
	glfw.KeyL:            `L`,
	glfw.KeyM:            `M`,
	glfw.KeyN:            `N`,
	glfw.KeyO:            `O`,
	glfw.KeyP:            `P`,
	glfw.KeyQ:            `Q`,
	glfw.KeyR:            `R`,
	glfw.KeyS:            `S`,
	glfw.KeyT:            `T`,
	glfw.KeyU:            `U`,
	glfw.KeyV:            `V`,
	glfw.KeyW:            `W`,
	glfw.KeyX:            `X`,
	glfw.KeyY:            `Y`,
	glfw.KeyZ:            `Z`,
	glfw.KeyLeftBracket:  `[`,
	glfw.KeyBackslash:    `\`,
	glfw.KeyRightBracket: `]`,
	glfw.KeyGraveAccent:  "`",
	glfw.KeyEscape:       `Esc`,
	glfw.KeyEnter:        `Enter`,
	glfw.KeyTab:          `Tab`,
	glfw.KeyBackspace:    `BackSpace`,
	glfw.KeyInsert:       `Insert`,
	glfw.KeyDelete:       `Delete`,
	glfw.KeyRight:        `Right`,
	glfw.KeyLeft:         `Left`,
	glfw.KeyDown:         `Down`,
	glfw.KeyUp:           `Up`,
	glfw.KeyPageUp:       `PageUp`,
	glfw.KeyPageDown:     `PageDown`,
	glfw.KeyHome:         `Home`,
	glfw.KeyEnd:          `End`,
	glfw.KeyCapsLock:     `CapsLock`,
	glfw.KeyScrollLock:   `ScrollLock`,
	glfw.KeyNumLock:      `NumLock`,
	glfw.KeyPrintScreen:  `PrintScreen`,
	glfw.KeyPause:        `Pause`,
	glfw.KeyF1:           `F1`,
	glfw.KeyF2:           `F2`,
	glfw.KeyF3:           `F3`,
	glfw.KeyF4:           `F4`,
	glfw.KeyF5:           `F5`,
	glfw.KeyF6:           `F6`,
	glfw.KeyF7:           `F7`,
	glfw.KeyF8:           `F8`,
	glfw.KeyF9:           `F9`,
	glfw.KeyF10:          `F10`,
	glfw.KeyF11:          `F11`,
	glfw.KeyF12:          `F12`,
	glfw.KeyF13:          `F13`,
	glfw.KeyF14:          `F14`,
	glfw.KeyF15:          `F15`,
	glfw.KeyF16:          `F16`,
	glfw.KeyF17:          `F17`,
	glfw.KeyF18:          `F18`,
	glfw.KeyF19:          `F19`,
	glfw.KeyF20:          `F20`,
	glfw.KeyF21:          `F21`,
	glfw.KeyF22:          `F22`,
	glfw.KeyF23:          `F23`,
	glfw.KeyF24:          `F24`,
	glfw.KeyF25:          `F25`,
	glfw.KeyKP0:          `0`,
	glfw.KeyKP1:          `1`,
	glfw.KeyKP2:          `2`,
	glfw.KeyKP3:          `3`,
	glfw.KeyKP4:          `4`,
	glfw.KeyKP5:          `5`,
	glfw.KeyKP6:          `6`,
	glfw.KeyKP7:          `7`,
	glfw.KeyKP8:          `8`,
	glfw.KeyKP9:          `9`,
	glfw.KeyKPDecimal:    `.`,
	glfw.KeyKPDivide:     `/`,
	glfw.KeyKPMultiply:   `*`,
	glfw.KeyKPSubtract:   `-`,
	glfw.KeyKPAdd:        `+`,
	glfw.KeyKPEnter:      `Enter`,
	glfw.KeyKPEqual:      `=`,
	glfw.KeyLeftShift:    `LShift`,
	glfw.KeyLeftControl:  `LCtrl`,
	glfw.KeyLeftAlt:      `LAlt`,
	glfw.KeyLeftSuper:    `LSuper`,
	glfw.KeyRightShift:   `RShift`,
	glfw.KeyRightControl: `RCtrl`,
	glfw.KeyRightAlt:     `RAlt`,
	glfw.KeyRightSuper:   `RSuper`,
	glfw.KeyMenu:         `Menu`,
}
