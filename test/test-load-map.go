package main

import(
  //"fmt"
	"AI/models"
	//"AI/astar"
	//"path/filepath"
	//"os"
	//"time"
	//"io/ioutil"
)

func main(){
  //fmt.Println();
  //tmxMapIns, tsx := models.InitMapStaticResource();
  tmx, _ := models.InitMapStaticResource();
  //tmxMapIns.PathFindingMap()


  //fmt.Println("1111111111111111111");
  //fmt.Println(tsx);

  //barriers := InitBarriers2(&tmx, &tsx);

  //fmt.Println(tsx);

  models.InitItemsForPathFinding(&tmx);
  models.FindPath(&tmx);


  //walkInfo := models.AstarPathToWalkInfo(tmxMapIns.Path);
  //step := 300.0;

  /*
  for {
    end := models.GotToGoal(step, &walkInfo);
    if end{
      break;
    }else{
     time.Sleep(1 * time.Second);
    }
  }
  */

}
