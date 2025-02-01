package prompter

import "fmt"

type FSMode struct {
	mode FSFlag
}

type FSFlag uint8

const (
	FileOnly FSFlag = 1 << iota
	DirectoryOnly
	FollowSymlinks
	MustExist
	MustNotExist
)

func (f FSFlag) String() string {
	switch f {
	case FileOnly:
		return "FileOnly"
	case DirectoryOnly:
		return "DirectoryOnly"
	case FollowSymlinks:
		return "FollowSymlinks"
	case MustExist:
		return "MustExist"
	case MustNotExist:
		return "MustNotExist"
	default:
		return "Unknown"
	}
}

func Debug(m FSMode) {
	fmt.Println("--- --- ---")
	fmt.Printf("Files allowed: %t\n", m.FileOnly())
	fmt.Printf("Directories allowed: %t\n", m.DirectoryOnly())
	fmt.Printf("Follow Symlinks: %t\n", m.FollowSymlinks())
	fmt.Printf("Must Exist: %t\n", m.MustExist())
}

func (m *FSMode) Set(flag FSFlag) {
	m.mode = m.mode | flag

	if flag == FileOnly {
		m.Set(MustExist)
		m.Clear(MustNotExist)
		m.Clear(DirectoryOnly)
	}

	if flag == DirectoryOnly {
		m.Set(MustExist)
		m.Clear(MustNotExist)
		m.Clear(FileOnly)
	}

	if flag == MustExist {
		m.Clear(MustNotExist)
	}

	if flag == MustNotExist {
		m.Clear(MustExist)
		m.Clear(FileOnly)
		m.Clear(DirectoryOnly)
	}
}

func (m *FSMode) Clear(flag FSFlag) {
	m.mode = m.mode &^ flag
}

func (m *FSMode) Toggle(flag FSFlag) {
	if m.Has(flag) {
		m.Clear(flag)
		return
	}

	m.Set(flag)
}

func (m *FSMode) Has(flag FSFlag) bool {
	return m.mode&flag != 0
}

func (m *FSMode) FileOnly() bool {
	return m.mode&FileOnly != 0
}

func (m *FSMode) DirectoryOnly() bool {
	return m.mode&DirectoryOnly != 0
}

func (m *FSMode) FollowSymlinks() bool {
	return m.mode&FollowSymlinks != 0
}

func (m *FSMode) MustExist() bool {
	return m.mode&MustExist != 0
}

func (m *FSMode) MustNotExist() bool {
	return m.mode&MustNotExist != 0
}
