package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	baseUrl    = "http://apis.data.go.kr/1230000/ao/PubDataOpnStdService/getDataSetOpnStdBidPblancInfo"
	serviceKey = ""
)

func createBidAnnounceRequest(pageNo, numOfRows int) (*http.Request, error) {
	params := url.Values{}
	params.Add("serviceKey", serviceKey)
	params.Add("type", "json")
	params.Add("pageNo", strconv.Itoa(pageNo))
	params.Add("numOfRows", strconv.Itoa(numOfRows))

	reqUrl := fmt.Sprintf("%s?%s", baseUrl, params.Encode())
	req, err := http.NewRequest("GET", reqUrl, nil)

	if err != nil {
		return nil, err
	}
	return req, nil
}

func FetchBidAnnouncements() ([]BidItem, error) {
	req, err := createBidAnnounceRequest(1, 10)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("not ok")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var apiResp BidAnnouncementAPIResponse
	if err = json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}
	return apiResp.Response.Body.Items, nil
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	bidAnnouncements, err := FetchBidAnnouncements()
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(bidAnnouncements)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(bytes)), nil
}

type BidAnnouncementAPIResponse struct {
	Response struct {
		Body struct {
			Items []BidItem
		}
	}
}

type BidItem struct {
	BidNtceNo                  string `json:"bidNtceNo"`                  // 입찰공고번호
	BidNtceOrd                 string `json:"bidNtceOrd"`                 // 입찰공고차수
	RefNtceNo                  string `json:"refNtceNo"`                  // 참조공고번호
	RefNtceOrd                 string `json:"refNtceOrd"`                 // 참조공고차수
	PpsNtceYn                  string `json:"ppsNtceYn"`                  // 사전공고여부
	BidNtceNm                  string `json:"bidNtceNm"`                  // 입찰공고명
	BidNtceSttusNm             string `json:"bidNtceSttusNm"`             // 입찰공고상태명
	BidNtceDate                string `json:"bidNtceDate"`                // 입찰공고일자
	BidNtceBgn                 string `json:"bidNtceBgn"`                 // 입찰공고시작시간
	BsnsDivNm                  string `json:"bsnsDivNm"`                  // 업무구분명
	IntrntnlBidYn              string `json:"intrntnlBidYn"`              // 국제입찰여부
	CmmnCntrctYn               string `json:"cmmnCntrctYn"`               // 공동계약여부
	CmmnReciptMethdNm          string `json:"cmmnReciptMethdNm"`          // 공동접수방식명
	ElctrnBidYn                string `json:"elctrnBidYn"`                // 전자입찰여부
	CntrctCnclsSttusNm         string `json:"cntrctCnclsSttusNm"`         // 계약체결상태명
	CntrctCnclsMthdNm          string `json:"cntrctCnclsMthdNm"`          // 계약체결방법명
	BidwinrDcsnMthdNm          string `json:"bidwinrDcsnMthdNm"`          // 낙찰자결정방법명
	NtceInsttNm                string `json:"ntceInsttNm"`                // 공고기관명
	NtceInsttCd                string `json:"ntceInsttCd"`                // 공고기관코드
	NtceInsttOfclDeptNm        string `json:"ntceInsttOfclDeptNm"`        // 공고기관담당부서명
	NtceInsttOfclNm            string `json:"ntceInsttOfclNm"`            // 공고기관담당자명
	NtceInsttOfclTel           string `json:"ntceInsttOfclTel"`           // 공고기관담당자전화번호
	NtceInsttOfclEmailAdrs     string `json:"ntceInsttOfclEmailAdrs"`     // 공고기관담당자이메일주소
	DmndInsttNm                string `json:"dmndInsttNm"`                // 수요기관명
	DmndInsttCd                string `json:"dmndInsttCd"`                // 수요기관코드
	DmndInsttOfclDeptNm        string `json:"dmndInsttOfclDeptNm"`        // 수요기관담당부서명
	DmndInsttOfclNm            string `json:"dmndInsttOfclNm"`            // 수요기관담당자명
	DmndInsttOfclTel           string `json:"dmndInsttOfclTel"`           // 수요기관담당자전화번호
	DmndInsttOfclEmailAdrs     string `json:"dmndInsttOfclEmailAdrs"`     // 수요기관담당자이메일주소
	PresnatnOprtnYn            string `json:"presnatnOprtnYn"`            // 제안서설명회여부
	PresnatnOprtnDate          string `json:"presnatnOprtnDate"`          // 제안서설명회일자
	PresnatnOprtnTm            string `json:"presnatnOprtnTm"`            // 제안서설명회시간
	PresnatnOprtnPlce          string `json:"presnatnOprtnPlce"`          // 제안서설명회장소
	BidPrtcptQlfctRgstClseDate string `json:"bidPrtcptQlfctRgstClseDate"` // 입찰참가자격등록마감일자
	BidPrtcptQlfctRgstClseTm   string `json:"bidPrtcptQlfctRgstClseTm"`   // 입찰참가자격등록마감시간
	CmmnReciptAgrmntClseDate   string `json:"cmmnReciptAgrmntClseDate"`   // 공동수급협정마감일자
	CmmnReciptAgrmntClseTm     string `json:"cmmnReciptAgrmntClseTm"`     // 공동수급협정마감시간
	BidBeginDate               string `json:"bidBeginDate"`               // 입찰시작일자
	BidBeginTm                 string `json:"bidBeginTm"`                 // 입찰시작시간
	BidClseDate                string `json:"bidClseDate"`                // 입찰마감일자
	BidClseTm                  string `json:"bidClseTm"`                  // 입찰마감시간
	OpengDate                  string `json:"opengDate"`                  // 개찰일자
	OpengTm                    string `json:"opengTm"`                    // 개찰시간
	OpengPlce                  string `json:"opengPlce"`                  // 개찰장소
	AsignBdgtAmt               string `json:"asignBdgtAmt"`               // 배정예산금액
	PresmptPrce                string `json:"presmptPrce"`                // 추정가격
	RsrvtnPrceDcsnMthdNm       string `json:"rsrvtnPrceDcsnMthdNm"`       // 예정가격결정방법명
	RgnLmtYn                   string `json:"rgnLmtYn"`                   // 지역제한여부
	PrtcptPsblRgnNm            string `json:"prtcptPsblRgnNm"`            // 참가가능지역명
	IndstrytyLmtYn             string `json:"indstrytyLmtYn"`             // 업종제한여부
	BidprcPsblIndstrytyNm      string `json:"bidprcPsblIndstrytyNm"`      // 입찰참가가능업종명
	BidNtceUrl                 string `json:"bidNtceUrl"`                 // 입찰공고URL
	DataBssDate                string `json:"dataBssDate"`                // 데이터기준일자
}

func main() {
	log.Println("server started")
	serviceKey = os.Getenv("serviceKey")
	s := server.NewMCPServer(
		"g2b",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	searchTool := mcp.NewTool("search_bid_announcements",
		mcp.WithDescription("나라장터 공고 검색 mcp 툴입니다."),
		mcp.WithString(
			"bidStartDate",
			mcp.Description("입찰 공고 시작일, YYYYMMDDHHMM"),
		),
		mcp.WithString(
			"bidEndDate",
			mcp.Description("입찰 공고 종료일, YYYYMMDDHHMM"),
		),
	)

	s.AddTool(searchTool, helloHandler)
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
