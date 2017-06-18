package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	//"time"
	"os"
	"os/exec"
	"bufio"
	"strings"
	"errors"
)

const COMMAND_END = "ENDNdgjTy2GpDGuzE4ueZeaRrSxHut"

type content struct {
	line string
	//filename string //barcode
	Select bool
}

type command struct{
	display string
	function string
	key rune
}

type display_and_contents struct{
     Display_Results []string
     Contents []content   
}

var commands []command = []command{
	{"[I] import", "import", 'i'},
	{"[E] export", "export", 'e'},
	{"[S] escape", "escape", 's'},
	{"[F] format", "format", 'f'},
}


func exec_import(ch chan string, barcodes ...string){
	for _, barcode := range barcodes{
		var _ = barcode
		cmd := exec.Command("sh", "echo.sh", barcode)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
		fmt.Println(err)
		os.Exit(1)
		}

		cmd.Start()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			//fmt.Println(scanner.Text())
			//fmt.Println()
			ch<-scanner.Text()
		}

		cmd.Wait()
		
	}
	ch<-COMMAND_END
	close(ch)
}

func line_to_barcode(line string) (string, error){
	array := strings.Fields(line)
	if len(array) < 8 { return  "", errors.New("not barcode")}
	return array[0], nil
}

func execute_batch(cmd command, contents ...content){
	barcodes := []string{}
	ch := make(chan string)
	for _, c := range contents{
		if(c.Select) {
			b, err := line_to_barcode(c.line)
			if err != nil {
			}
			barcodes = append(barcodes, b) 
		}
	}
	if cmd.function == "import" {
		//fmt.Println("debug")
		go exec_import(ch, barcodes...)
	} else if cmd.function == "export" {

	} else if cmd.function == "escape" {

	} else if cmd.function == "format" {

	}
	display_result := []string{}
	for {
		result := <-ch
		if(result == COMMAND_END){ break }
		display_result = append(display_result, result)
	    a := display_and_contents{Display_Results:display_result,Contents:contents}
	    result_draw(a)
		}

}

func result_draw(a display_and_contents) {
	defer termbox.Flush()
	results := a.Display_Results
	contents := a.Contents
	w, h := termbox.Size()
	draw(0, contents...)
	for i, r := range results{
	tbPrint(w/2, h/2 + i, termbox.ColorRed, termbox.ColorDefault, string(r))
	}
}

func tbPrint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func ask_draw(ask_cur bool, c command, cursor int, contents ...content){
	defer termbox.Flush()
	_, h := termbox.Size()
	var caution string = "Is " + string(c.display) + " operation is OK?"
	draw(cursor, contents...)
	tbPrint(3, h - 2, termbox.ColorRed, termbox.ColorDefault, caution)
	if(ask_cur){
	     tbPrint(3 + 40, h - 2, termbox.ColorRed, termbox.ColorDefault | termbox.AttrReverse, "[  YES  ]")
	     tbPrint(3 + 50, h - 2, termbox.ColorRed, termbox.ColorDefault, "[  NO  ]")
	     }else
	     {
	     tbPrint(3 + 40, h - 2, termbox.ColorRed, termbox.ColorDefault, "[  YES  ]")
	     tbPrint(3 + 50, h - 2, termbox.ColorRed, termbox.ColorDefault | termbox.AttrReverse, "[  NO  ]")
	 }
}

func ask_before_exec(c command, cursor int, contents ...content) bool {
	
	ask_cur := false
	ask_draw(ask_cur, c, cursor, contents...)
askloop:
	for
	{
		switch ev:= termbox.PollEvent(); ev.Type{
		case termbox.EventInterrupt:
			break askloop
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlC {
				break askloop
			}else if ev.Key == termbox.KeyArrowLeft {
				ask_cur = true
			}else if ev.Key == termbox.KeyArrowRight {
				ask_cur = false
			}else if ev.Key == 0x0D {
				return ask_cur
			}

		}
		ask_draw(ask_cur, c, cursor, contents...)
	}
	return false
}


func draw(cursor int, contents ...content){
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()

	w, h := termbox.Size()
	s := fmt.Sprintf("hello world")
	tbPrint((w/2)-(len(s)/2), h/10, termbox.ColorRed, termbox.ColorDefault, s)
	for i, c := range contents{
		s := fmt.Sprintf(string(c.line))
		fgattr := termbox.ColorRed
		bgattr := termbox.ColorDefault
		if(c.Select){
			bgattr = bgattr | termbox.AttrReverse
		} else {
			bgattr = termbox.ColorDefault
		}
		if(i == cursor) {
			fgattr = fgattr | termbox.AttrUnderline
		}

		tbPrint(/*(w/2)-(len(s)/2)*/ h/10, i + 1 + h/10, fgattr, bgattr, s)
	}

    /* command display */

    for i, c  := range commands{
    	s := fmt.Sprintf(string(c.display))
    	fgattr := termbox.ColorRed
		bgattr := termbox.ColorDefault
    	tbPrint(12 * i + 1, h - 1, fgattr, bgattr, s)
    }

}

func key_control(ev *termbox.Event){

}

func getInformation()([]content){
	contents := []content{}
	cmd := exec.Command("ls", "-la")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		//fmt.Println()
		c := content{line: scanner.Text()} 
		contents = append(contents, c)
	}
	cmd.Wait()

	return contents
}

func check_content(contents *[]content){

}





func main(){
	
	contents := getInformation()
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	/*go func(){
		time.Sleep(5 * time.Second)
		termbox.Interrupt()
		//time.Sleep(1 * time.Second)
		//panic("this should never run")
	}()
	*/
	
	cursor := 0
	draw(cursor, contents...)
mainloop:
	for
	{
    	switch ev:= termbox.PollEvent(); ev.Type{
    	case termbox.EventInterrupt:
    		break mainloop
    	case termbox.EventKey:
    		if ev.Key == termbox.KeyCtrlC {
    			break mainloop
    		}else if ev.Key == termbox.KeyArrowDown {
    			cursor++
    			if(cursor > len(contents) - 1){ cursor = len(contents) - 1 }
    			break
    		}else if ev.Key == termbox.KeyArrowUp {
    			cursor--
    			if(cursor < 0){ cursor = 0 }
    			break
    		} else if ev.Key == 0x0D {/* Enter */
    			if(contents[cursor].Select){
    				contents[cursor].Select = false
    			}else{
    			    contents[cursor].Select = true
    			}
    			break
    		}
    		for _, c := range commands{
    			if(c.key == ev.Ch){
    				ans := ask_before_exec(c, cursor, contents...)
    				if(ans == true){
    					execute_batch(c, contents...)
    					draw(cursor, contents...)
    				}
    				break
    			}
    		}

    	}
    	draw(cursor, contents...)
    }

	termbox.Close()
	fmt.Println("Finished")
	os.Exit(0)
}










