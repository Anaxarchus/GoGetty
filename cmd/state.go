package cmd

import (
	"fmt"
	"gogetty/pkg/app"
	"gogetty/pkg/cache"
	"gogetty/pkg/gitop"
	"os"
)

func getApp() *app.MyApp {
	var myApp *app.MyApp
	workingDir, wdErr := os.Getwd()
	if wdErr != nil {
		fmt.Println("Error getting working directory:", wdErr)
		return nil
	}
	modules, err := gitop.Scan(cache.ModuleDir())
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	myApp = &app.MyApp{
		ProjectDir: workingDir,
		Cache:      modules,
	}
	return myApp
}
