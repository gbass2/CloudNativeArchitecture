package gameapi

import (
  "testing"
  "encoding/json"
  "fmt"
)

func TestSearchGame(t *testing.T) {
  data := make(map[string]interface{})
  data["name"] = "Halo: Combat Evolved"

  var want []map[string]interface{} // Slice of map of the returned data.
  want = append(want, data)

  got, err := SearchGame("Halo: Combat Evolved")

  if err != nil {
    t.Error("Error getting game information.")
  }

  gotJson, _ := json.MarshalIndent(got, "", "    ")
  fmt.Println(string(gotJson))


  if want[0]["name"] != got[0]["name"] {
      t.Error("Error in api.gameSearch; \n\n Want " + want[0]["name"].(string) + "\n Got " + got[0]["name"].(string))
  }
}
