package endpoints

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/connectionManager/endpoints/internal"
	"github.com/reactivex/rxgo/v2"
	"html/template"
	"net/http"
	"sort"
)

func GetConnections(Template *template.Template, ConnectionManager connectionManager.IObtainConnectionManagerInformation) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		connections, err := ConnectionManager.GetConnections(context.TODO())
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		sort.Sort(&internal.SortConnectionInformation{Connections: connections})
		err = Template.Execute(writer, &struct {
			Connections []*internal.ConnectionInformation
		}{
			Connections: func(data []*connectionManager.ConnectionInformation) []*internal.ConnectionInformation {
				var result []*internal.ConnectionInformation
				for _, ci := range data {

					StackPropertiesMap := make(map[int]*internal.StackData)
					for key, value := range ci.StackProperties {
						dd, ok := StackPropertiesMap[key.Index]
						if !ok {
							dd = &internal.StackData{
								Name:  key.Name,
								Index: key.Index,
							}
							StackPropertiesMap[key.Index] = dd
						}
						switch key.Direction {
						case rxgo.StreamDirectionInbound:
							dd.InValue = value
						case rxgo.StreamDirectionOutbound:
							dd.OutValue = value
						}
						continue

					}
					var StackPropertiesArray []*internal.StackData
					for _, value := range StackPropertiesMap {
						StackPropertiesArray = append(StackPropertiesArray, value)
					}
					sort.Sort(&internal.StackDataSort{Data: StackPropertiesArray})

					result = append(result, &internal.ConnectionInformation{
						Id:              ci.Id,
						Name:            ci.Name,
						Status:          ci.Status,
						ConnectionTime:  ci.ConnectionTime,
						StackProperties: StackPropertiesArray,
					})
				}
				return result
			}(connections),
		})
	}
}
