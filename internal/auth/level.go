package auth

type Level uint8

const (
	Unlogged Level = iota

	valid_start

	Chief
	Promoted
	Admin
	User

	valid_end
)

func (l Level) IsValid() bool {
	return valid_start < l && l < valid_end
}

func (l Level) IsValidOrUnlogged() bool {
	return l == Unlogged || valid_start < l && l < valid_end
}

func (l Level) String() string {
	return levelStrings[l]
}

func FromString(string string) (Level, bool) {
	level, in := stringLevels[string]
	if !in {
		return Unlogged, false
	}

	return level, true
}

var levelStrings = map[Level]string{
	Chief:    "chief",
	Promoted: "promoted-admin",
	Admin:    "admin",
	User:     "user",

	Unlogged: "unlogged",
}

var stringLevels = mirror(levelStrings)

func mirror[K, V comparable](m map[K]V) map[V]K {
	nm := make(map[V]K, len(m))
	for k, v := range m {
		nm[v] = k
	}
	return nm
}
