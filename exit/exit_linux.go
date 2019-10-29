package exit

func init() {
	signalType = initSignal()
	signalType = append(signalType, syscall.SIGUSR1, syscall.SIGUSR2)
}