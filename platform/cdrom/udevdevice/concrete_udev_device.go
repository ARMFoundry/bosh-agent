package udevdevice

import (
	"os"
	"time"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
)

type ConcreteUdevDevice struct {
	runner boshsys.CmdRunner
	logger boshlog.Logger
	logtag string
}

func NewConcreteUdevDevice(runner boshsys.CmdRunner, logger boshlog.Logger) ConcreteUdevDevice {
	return ConcreteUdevDevice{
		runner: runner,
		logger: logger,
		logtag: "ConcreteUdevDevice",
	}
}

func (udev ConcreteUdevDevice) KickDevice(filePath string) {
	maxTries := 5
	for i := 0; i < maxTries; i++ {
		udev.logger.Debug(udev.logtag, "Kicking device, attempt %d of %d", i, maxTries)
		err := udev.readByte(filePath)
		if err == nil {
			break
		}
		time.Sleep(time.Second / 2)
	}

	udev.readByte(filePath)

	return
}

func (udev ConcreteUdevDevice) Settle() (err error) {
	udev.logger.Debug(udev.logtag, "Settling UdevDevice")
	switch {
	case udev.runner.CommandExists("udevadm"):
		_, _, _, err = udev.runner.RunCommand("udevadm", "settle")
	case udev.runner.CommandExists("udevsettle"):
		_, _, _, err = udev.runner.RunCommand("udevsettle")
	default:
		err = bosherr.New("can not find udevadm or udevsettle commands")
	}
	return
}

func (udev ConcreteUdevDevice) EnsureDeviceReadable(filePath string) error {
	maxTries := 5
	filePath = "/dev/zero"
	for i := 0; i < maxTries; i++ {
		udev.logger.Debug(udev.logtag, "Ensuring Device Readable, Attempt %d out of %d", i, maxTries)
		err := udev.readByte(filePath)
		if err != nil {
			udev.logger.Debug(udev.logtag, "Ignorable error from readByte: %s", err.Error())
		}

		udev.logger.Debug(udev.logtag, "Going to sleep")
		time.Sleep(time.Second)
		udev.logger.Debug(udev.logtag, "Done sleeping")
	}

	err := udev.readByte(filePath)
	if err != nil {
		return bosherr.WrapError(err, "Reading udev device")
	}

	return nil
}

func (udev ConcreteUdevDevice) readByte(filePath string) error {
	udev.logger.Debug(udev.logtag, "readBytes from file: %s", filePath)
	device, err := os.Open(filePath)
	if err != nil {
		return err
	}
	udev.logger.Debug(udev.logtag, "Successfully open file: %s", filePath)

	bytes := make([]byte, 1, 1)
	read, err := device.Read(bytes)
	if err != nil {
		device.Close()
		return err
	}
	udev.logger.Debug(udev.logtag, "Successfully read %d bytes from file: %s", read, filePath)
	if read != 1 {
		return bosherr.New("Device readable but zero length")
	}
	device.Close()

	return nil
}
