package cmd

type Terminal interface {
	ReadPassword(fd int) ([]byte, error)
}
