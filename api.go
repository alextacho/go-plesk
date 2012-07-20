package plesk

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

func NewPleskApi(url string, username string, password string) *PleskAPI {
	return &PleskAPI{
		url: url,
		username: username,
		password: password,
	}
}

type PleskAPI struct {
	url      string
	username string
	password string
}

type getMailResponse struct {
	Status             string   `xml:"mail>get_info>result>status"`
	Errortxt           string   `xml:"mail>get_info>result>errtext"`
	ForwardingAdresses []string `xml:"mail>get_info>result>mailname>forwarding>address"`
}

type getMailingListResponse struct {
	Status   string `xml:"maillist>get-list>result>status"`
	Errortxt string `xml:"maillist>get-list>result>errtext"`
}

func (self *PleskAPI) CreateEmail(name string, forward []string) {
	// POST requets
	// debug.Nop()

	packet := `<packet version='1.6.3.0'>
				    <mail>
				   	    <create>
				            <filter>
				                <site-id>55</site-id>
				                <mailname>
				                    <name>` + name + `</name>
			                        <mailbox>
              							<enabled>false</enabled>
              						</mailbox>
		                            <forwarding>
					                    <enabled>true</enabled>`
	for i := 0; i < len(forward); i++ {
		packet += `<address>` + forward[i] + `</address>`
	}
	packet += `</forwarding>
	                            </mailname> 
	                        </filter>
						</create>
					</mail>
				</packet>`

	self.doRequest(packet)

}

func (self *PleskAPI) CreateMailingList(list string) {
	// POST requets
	// debug.Nop()

	packet := `<packet version='1.6.3.0'>
				    <maillist>
						<add-list>
						   <site-id>55</site-id>
						   <name>` + list + `</name>
						   <password>stopeurope</password>
						   <admin-email>it@starteurope.at</admin-email>
						</add-list>
					</maillist>
				</packet>`

	self.doRequest(packet)

}

func (self *PleskAPI) UpdateEmail(name string, forward []string) {
	// POST requets
	// debug.Nop()

	packet := `<packet version='1.6.3.0'>
				    <mail>
				   	    <update>
				   	    	<set>
					            <filter>
					                <site-id>55</site-id>
					                <mailname>
					                    <name>` + name + `</name>
			                            <forwarding>
						                    <enabled>true</enabled>`
	for i := 0; i < len(forward); i++ {
		packet += `<address>` + forward[i] + `</address>`
	}

	packet += `			                </forwarding>
		                            </mailname> 
		                        </filter>
	                        </set>
						</update>
					</mail>
				</packet>`

	self.doRequest(packet)

}

func (self *PleskAPI) AddEmailToList(plesk PleskAPI, list string, name string) {
	// POST requets
	// debug.Nop()

	packet := `<packet version='1.6.3.0'>
				    <maillist>
						<add-member>
						   <filter>
						      <list-name>` + list + `</list-name>						    
						   </filter>
						   <id>` + name + `</id>
						</add-member>
					</maillist>
				</packet>`

	self.doRequest(packet)

}

func (self *PleskAPI) EmailExists(name string) (resp getMailResponse, err error) {
	// POST requets
	// debug.Nop()
	var result getMailResponse

	packet := `<packet version="1.6.3.0">
					<mail>
				  		<get_info>
				   			<filter>
				   				<site-id>55</site-id>
              					<name>` + name + `</name>
				   			</filter>
				   			<forwarding></forwarding>
				  		</get_info>
				 	</mail>
				</packet>`

	// reader := strings.Read([]byte(packet))

	response, err := self.doRequest(packet)

	if err != nil {
		return result, err
	}

	err = xml.Unmarshal(response, &result)
	if err != nil {
		// debug.Print("error: ", err)
		return result, err
	}

	return result, nil
}

func (self *PleskAPI) MailingListExists(listname string) (resp getMailingListResponse, err error) {
	// POST requets
	// debug.Nop()
	var result getMailingListResponse

	packet := `<packet version="1.6.3.0">
					<maillist>
						<get-list>
				   			<filter>
              					<name>` + listname + `</name>
				   			</filter>
				  		</get-list>
				 	</maillist>
				</packet>`

	// reader := strings.Read([]byte(packet))

	response, err := self.doRequest(packet)

	if err != nil {
		return result, err
	}

	err = xml.Unmarshal(response, &result)
	if err != nil {
		// debug.Print("error: ", err)
		return result, err
	}

	return result, nil
}

func (self *PleskAPI) doRequest(packet string) (response []byte, err error) {
	// POST requets
	// debug.Nop()
	url := self.url

	b := bytes.NewBufferString(packet)

	// debug.Print("body: ", b)
	// debug.Print("URL: " + url)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: tr,
	}

	req, err := http.NewRequest("POST", url, b)
	// ...
	req.Header.Add("HTTP_AUTH_LOGIN", self.username)
	req.Header.Add("HTTP_AUTH_PASSWD", self.password)
	req.Header.Add("Content-Type", "text/xml")

	resp, err := client.Do(req)

	// debug.Print("Response: ", resp)
	// debug.Print("Status: ", err)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	j, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// result := string(j)
	// debug.Print("Result:", result)

	return j, nil
}
