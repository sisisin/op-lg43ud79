package oplg43ud79

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/sisisin/op-lg43ud79/pkg/lg43client"
)

const (
	VID   = "0403"
	PID   = "6001"
	SetId = "01"
)

func RunSwitchDP(ctx context.Context, osArgs []string) error {
	f := flag.NewFlagSet("switchdp", flag.ExitOnError)
	debug := f.Bool("debug", false, "enable debug logging")
	f.Parse(osArgs[1:])

	if *debug {
		ctx = withLogLevel(ctx, LogLevelDebug)
		ctx = lg43client.WithLogLevel(ctx, lg43client.LogLevelDebug)
	}

	logDebug(ctx, "instantiating LG43Client with PID=%s VID=%s", PID, VID)
	lg43, err := lg43client.NewLG43Client(ctx, lg43client.VID(VID), lg43client.PID(PID))
	if err != nil {
		return errors.Wrap(err, "NewLG43Client failed")
	}
	defer lg43.Close()

	logInfo(ctx, "prepare to switch to DP1")
	// NOTE: wake up display from sleep
	if res, err := lg43.PowerOn(ctx, SetId); err != nil {
		return errors.Wrap(err, "PowerOn failed")
	} else {
		logInfo(ctx, "PowerOn response: %s", res)
	}

	logInfo(ctx, "switching to DP1")
	if res, err := lg43.InputSelectToDP1(ctx, SetId); err != nil {
		return errors.Wrap(err, "InputSelectToDP1 failed")
	} else {
		logInfo(ctx, "InputSelectToDP1 response: %s", res)
	}

	return nil
}

func RunSwitchHDMI4(ctx context.Context, osArgs []string) error {
	f := flag.NewFlagSet("switchhdmi4", flag.ExitOnError)
	debug := f.Bool("debug", false, "enable debug logging")
	f.Parse(osArgs[1:])

	if *debug {
		ctx = withLogLevel(ctx, LogLevelDebug)
		ctx = lg43client.WithLogLevel(ctx, lg43client.LogLevelDebug)
	}
	lg43, err := lg43client.NewLG43Client(ctx, lg43client.VID(VID), lg43client.PID(PID))
	if err != nil {
		return errors.Wrap(err, "NewLG43Client failed")
	}
	defer lg43.Close()

	// NOTE: wake up display from sleep
	if res, err := lg43.PowerOn(ctx, SetId); err != nil {
		return errors.Wrap(err, "PowerOn failed")
	} else {
		logInfo(ctx, "PowerOn response: %s", res)
	}

	if res, err := lg43.InputSelectToHDMI4(ctx, SetId); err != nil {
		return errors.Wrap(err, "InputSelectToHDMI4 failed")
	} else {
		logInfo(ctx, "InputSelectToHDMI4 response: %s", res)
	}

	return nil
}

func RunWriteLG43(ctx context.Context, osArgs []string) error {
	f := flag.NewFlagSet("writelg43", flag.ExitOnError)
	debug := f.Bool("debug", false, "enable debug logging")
	f.Parse(osArgs[1:])

	if *debug {
		ctx = withLogLevel(ctx, LogLevelDebug)
		ctx = lg43client.WithLogLevel(ctx, lg43client.LogLevelDebug)
	}

	args := f.Args()
	if len(args) != 3 {
		fmt.Printf("error: invalid arguments\nreceived: %v\n", args)
		fmt.Println("Usage: op-lg43ud79 <command> <setId> <data>")
		fmt.Println("Example: op-lg43ud79 kh 01 1E # brightness to 30")
		os.Exit(1)
	}
	cmd, setId, data := args[0], args[1], args[2]

	lg43, err := lg43client.NewLG43Client(ctx, lg43client.VID(VID), lg43client.PID(PID))
	if err != nil {
		return errors.Wrap(err, "NewLG43Client failed")
	}
	defer lg43.Close()

	w := fmt.Sprintf("%s %s %s\r", cmd, setId, data)
	res, err := lg43.Write(ctx, []byte(w))
	if err != nil {
		return errors.Wrap(err, "Write failed")
	}
	logInfo(ctx, "response: %s", res)
	return nil
}
