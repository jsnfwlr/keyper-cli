package feedback

var cutOff Level = 0

func SetVerbosity(levels Level) {
	switch {
	case levels > LevelCap:
		levels = LevelCap
	case levels < 0:
		levels = 0
	}

	cutOff = levels
}

func SuppressNoise(levels Level) {
	target := (LevelCap - levels)
	if target < 0 {
		target = 0
	}
	cutOff = target
}

func GetCutOff() Level {
	return cutOff
}

var beforeMute Level

func Mute() {
	beforeMute = GetCutOff()
	SuppressNoise(LevelCap)
}

func Unmute() {
	SuppressNoise(NoiseLevel(int(LevelCap) - int(beforeMute)))
}

func ExceedsLimit(level Level) bool {
	return level > cutOff
}

func NoiseLevel(level int) Level {
	return Level(level)
}
