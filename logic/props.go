package logic

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/levigross/grequests"
	"github.com/pkg/errors"
)

func Fetch_props(cid int, token string) error {
	token = strings.TrimSpace(token)
	if cid < 0 || token == "" {
		return errors.New("invalid cid or token")
	}

	url := fmt.Sprintf("https://open.taobao.com/handler/tools/getApiResult.json?isEn=false&env=2&format=json&apiName=taobao.itemprops.get&appKey=&fields=&cid=%d&pid=&parent_pid=&is_key_prop=&is_sale_prop=&is_color_prop=&is_enum_prop=&is_input_prop=&is_item_prop=&child_path=&type=&attr_keys=&appSource=0&_tb_token_=%s", cid, token)
	resp, err := grequests.Get(url, &grequests.RequestOptions{
		Cookies: []*http.Cookie{
			&http.Cookie{
				Name:  "_tb_token_",
				Value: token,
			},
		},
	})
	if err != nil {
		return errors.WithStack(err)
	}
	if !resp.Ok {
		return errors.New("resp StatusCode isn't 2xx")
	}

	var data PropsData
	if err = resp.JSON(&data); err != nil {
		return errors.WithStack(err)
	}

	var propsResp PropsResponse
	if err = json.Unmarshal([]byte(data.Data.Response), &propsResp); err != nil {
		return errors.WithStack(err)
	}

	if len(propsResp.ItempropsGetResponse.ItemProps.ItemProp) == 0 {
		return errors.New("该类目淘宝返回结果为空")
	}

	file, err := os.Create(fmt.Sprintf("%d.csv", cid))
	if err != nil {
		return errors.WithStack(err)
	}

	writer := csv.NewWriter(file)
	if err = writer.Write([]string{
		"类目id",
		"类目名",
		"属性id",
		"属性名",
		"属性值id",
		"属性值",
		"......",
	}); err != nil {
		return errors.WithStack(err)
	}

	for idx, prop := range propsResp.ItempropsGetResponse.ItemProps.ItemProp {
		values := prop.PropValues.PropValue
		row := make([]string, 0, len(values)*2+3)
		if idx == 0 {
			row = append(row, fmt.Sprintf("%d", cid), "类目名")
		} else {
			row = append(row, "", "")
		}

		row = append(row, fmt.Sprintf("%d", prop.Pid), prop.Name)

		for _, value := range values {
			row = append(row, fmt.Sprintf("%d", value.Vid), value.Name)
		}

		if err = writer.Write(row); err != nil {
			return errors.WithStack(err)
		}
	}

	writer.Flush()
	return nil
}

type PropsResponse struct {
	ItempropsGetResponse struct {
		LastModified string `json:"last_modified"`
		ItemProps    struct {
			ItemProp []struct {
				Pid         int    `json:"pid"`
				ParentPid   int    `json:"parent_pid"`
				ParentVid   int    `json:"parent_vid"`
				Name        string `json:"name"`
				IsKeyProp   bool   `json:"is_key_prop"`
				IsSaleProp  bool   `json:"is_sale_prop"`
				IsColorProp bool   `json:"is_color_prop"`
				IsEnumProp  bool   `json:"is_enum_prop"`
				IsItemProp  bool   `json:"is_item_prop"`
				Must        bool   `json:"must"`
				Multi       bool   `json:"multi"`
				PropValues  struct {
					PropValue []struct {
						Cid       int    `json:"cid"`
						Pid       int    `json:"pid"`
						PropName  string `json:"prop_name"`
						Vid       int    `json:"vid"`
						Name      string `json:"name"`
						NameAlias string `json:"name_alias"`
						IsParent  bool   `json:"is_parent"`
						Status    string `json:"status"`
						SortOrder int    `json:"sort_order"`
						Features  struct {
							Feature []struct {
								AttrKey   string `json:"attr_key"`
								AttrValue string `json:"attr_value"`
							} `json:"feature"`
						} `json:"features"`
					} `json:"prop_value"`
				} `json:"prop_values"`
				Status        string `json:"status"`
				SortOrder     int    `json:"sort_order"`
				ChildTemplate string `json:"child_template"`
				IsAllowAlias  bool   `json:"is_allow_alias"`
				IsInputProp   bool   `json:"is_input_prop"`
				Features      struct {
					Feature []struct {
						AttrKey   string `json:"attr_key"`
						AttrValue string `json:"attr_value"`
					} `json:"feature"`
				} `json:"features"`
				IsTaosir bool `json:"is_taosir"`
				TaosirDo struct {
					StdUnitList struct {
						StdUnit []struct {
							AttrKey   string `json:"attr_key"`
							AttrValue string `json:"attr_value"`
						} `json:"std_unit"`
					} `json:"std_unit_list"`
					ExprElList struct {
						ExprEl []struct {
							Type        int    `json:"type"`
							Text        string `json:"text"`
							IsShowLabel bool   `json:"is_show_label"`
							IsLabel     bool   `json:"is_label"`
							IsInput     bool   `json:"is_input"`
						} `json:"expr_el"`
					} `json:"expr_el_list"`
					Type      int `json:"type"`
					Precision int `json:"precision"`
				} `json:"taosir_do"`
				IsMaterial bool `json:"is_material"`
				MaterialDo struct {
					Materials struct {
						ItemMateriaValueDO []struct {
							Name              string `json:"name"`
							NeedContentNumber bool   `json:"need_content_number"`
						} `json:"item_materia_value_d_o"`
					} `json:"materials"`
				} `json:"material_do"`
			} `json:"item_prop"`
		} `json:"item_props"`
	} `json:"itemprops_get_response"`
}

type PropsData struct {
	Data struct {
		Response string `json:"response"`
	} `json:"data"`
}
