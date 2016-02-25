package main

import (
	"github.com/noobiedog/maltegolocal"
	"fmt"
	"os"
)

func main() {


	lt := maltegolocal.ParseLocalArguments(os.Args)
	Domain := lt.Value

	TRX := maltegolocal.MaltegoTransform{}

	NewEnt := TRX.AddEntity("maltego.Domain", "Hello" + Domain)
	NewEnt.SetType("maltego.Domain")
	NewEnt.SetValue(Domain)
	NewEnt.AddDisplayInformation("<h3>Heading</h3><p>content here about" + Domain + "!</p>", "Other")
	NewEnt.AddProperty("Display Value", Domain, "nostrict", "True")
	NewEnt.SetLinkColor("#FF0000")
	NewEnt.SetWeight(200) 
	NewEnt.SetNote("Domain is " + Domain)

	TRX.AddUIMessage("completed!","Inform")
	
 	fmt.Println(TRX.ReturnOutput())
}