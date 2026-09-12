package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var (
	timeStamp = make(map[string]time.Time)
	only string
)


func watchInnerFile(watcher *fsnotify.Watcher,cwd string){
			err := watcher.Add(cwd)
		if err != nil {
			log.Fatal(err)
		}

			info , err := os.Stat(cwd)
			if err != nil{
				return
			}

			if info.IsDir(){
				entries , err:=os.ReadDir(cwd)
				if err != nil{
					return
				}

				for _,entry := range entries{
					watchInnerFile(watcher,filepath.Join(cwd,entry.Name()))
				}
			}
	}

var root = &cobra.Command{
	Use: "watchon",
	Version: "1.0.0",
	Short: "Watchserver automatically watch a server. it works like nodemon e.g watchserver go run . , watchserver npm run dev",
	Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("No command specified")
				return
			}
		watcher , err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}

		var cwd = "."
		if only != ""{
			cwd = only
		}

		defer watcher.Close()
		watchInnerFile(
			watcher,cwd,
		)

		

		fmt.Println("Watching")

		var runningExec *exec.Cmd

		for {
			select{
			case events :=<- watcher.Events:

				key := events.Name + ":" + events.Op.String()
				lastTime,exists := timeStamp[key]
				if exists && time.Since(lastTime)< 500*time.Millisecond{
					continue
				}



					

				if events.Op&fsnotify.Create == fsnotify.Create{
							if _ , err:= os.Stat(filepath.Join(cwd,events.Name));err != nil {
						continue
					}
									fmt.Println("Created",events.Name)

					
					err = watcher.Add(filepath.Join(cwd,events.Name))
					if err != nil {
						fmt.Println("Error watching",events.Name,err)
						continue
					}

				}

				if events.Op&fsnotify.Write == fsnotify.Write{
					fmt.Println("Modified:", events.Name)
				}

				if events.Op.String() == "DELETE" {
					fmt.Println("Removed:", events.Name)
					watcher.Remove(filepath.Join(cwd,events.Name))
				}

				if events.Op.String() == "RENAME" {

					fmt.Println("Renamed:", events.Name)

					
					watcher.Remove(filepath.Join(cwd,events.Name))
				}
											timeStamp[events.String()+events.Name] = time.Now()

				if runningExec != nil && runningExec.Process != nil{
					err = runningExec.Process.Kill()
					if err != nil {
						println("Couldn't kill")
						continue
					}
					runningExec.Wait()
				}
				runningExec = exec.Command(args[0],args[1:]...)
				println("Restarting server..")
				runningExec.Stdout = os.Stdout
				runningExec.Stderr = os.Stderr

				err := runningExec.Start()
				if err != nil {
					fmt.Println("start error:", err)
				}




			case err:= <-watcher.Errors:
				fmt.Println(err)

			}
		}
	},
}

func init(){
	root.Flags().StringVar(&only,"only","","Watch only a folder or file")
}

func Execute(){
	if err := root.Execute(); err != nil{
		log.Fatal(err)
	}
}