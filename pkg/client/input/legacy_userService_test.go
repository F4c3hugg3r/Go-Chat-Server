package input

var (
	name     = "Arndt"
	clientId = "clientId-DyGWNnLrLWnbuhf-LgBUAdAxdZf-U1pgRw"
	//clientId2  = "clientId2-yGWNnLrLWnbuhf-LgBUAdAxdZf-U1pgRw"
	authToken = "authId-5EDyGWNnLrLWnbuhf-LgBUAdAxdZf-U1pgRwc7ex1dt5EDyGWNnLrLWnbuhf-LgBUAdAxdZf-U1pgRw"
	//authToken2 = "authId2-EDyGWNnLrLWnbuhf-LgBUAdAxdZf-U1pgRwc7ex1dt5EDyGWNnLrLWnbuhf-LgBUAdAxdZf-U1pgRw"
)

// func TestRegister(t *testing.T) {
// 	clientService := &Client{
// 		clientId:   clientId,
// 		authToken:  authToken,
// 		reader:     bufio.NewReader(os.Stdin),
// 		writer:     io.Writer(os.Stdout),
// 		HttpClient: &http.Client{},
// 	}

// 	rsp := Response{Name: "authToken", Content: authToken}
// 	jsonRsp, _ := json.Marshal(rsp)

// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write(jsonRsp)
// 	}))
// 	defer ts.Close()
// 	clientService.authToken = ""

// 	err := clientService.Register(ts.URL)
// 	assert.Error(t, err)

// 	ts.URL = fmt.Sprintf("%s/users/%s", ts.URL, clientId)

// 	clientService.reader = bufio.NewReader(strings.NewReader(name + "\n"))
// 	err = clientService.Register(ts.URL)
// 	assert.NoError(t, err)
// 	assert.Equal(t, clientService.authToken, authToken)
// }

// func TestClient_ReceiveMessages(t *testing.T) {
// 	c := &Client{
// 		clientId:   clientId,
// 		authToken:  authToken,
// 		reader:     bufio.NewReader(os.Stdin),
// 		writer:     io.Writer(os.Stdout),
// 		HttpClient: &http.Client{},
// 	}

// 	type fields struct {
// 		clientName string
// 		clientId   string
// 		reader     *bufio.Reader
// 		writer     io.Writer
// 		authToken  string
// 		HttpClient *http.Client
// 	}
// 	type args struct {
// 		rsp    *Response
// 		log    string
// 		stdOut string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{
// 			name: "fine",
// 			args: args{
// 				rsp:    &Response{Name: "Max", Content: "wubbalubbadubdub"},
// 				stdOut: "Max: wubbalubbadubdub\n",
// 				log:    "",
// 			},
// 		},
// 		{
// 			name: "inactive",
// 			args: args{
// 				rsp:    &Response{Name: "inactive", Content: "wubbalubbadubdub"},
// 				stdOut: "",
// 				log:    "inactivity",
// 			},
// 		},
// 		{
// 			name: "content empty",
// 			args: args{
// 				rsp:    &Response{Name: "Max", Content: ""},
// 				stdOut: "",
// 				log:    "",
// 			},
// 		},
// 		{
// 			args: args{
// 				rsp:    &Response{Name: "Max", Content: "[{}]"},
// 				stdOut: "|\n+\n|\n",
// 				log:    "",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var output bytes.Buffer
// 			c.writer = &output

// 			var logOutput bytes.Buffer
// 			log.SetOutput(&logOutput)
// 			defer log.SetOutput(os.Stderr)

// 			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 				jsonRsp, _ := json.Marshal(tt.args.rsp)
// 				w.Write(jsonRsp)
// 			}))

// 			_, cancel := context.WithCancel(context.Background())
// 			c.ReceiveMessages(ts.URL, cancel)

// 			assert.Equal(t, tt.args.stdOut, output.String())
// 			assert.Contains(t, logOutput.String(), tt.args.log)
// 		})
// 	}
// }

// func TestClient_SendMessages(t *testing.T) {
// 	c := &Client{
// 		clientId:   clientId,
// 		authToken:  authToken,
// 		reader:     bufio.NewReader(os.Stdin),
// 		writer:     io.Writer(os.Stdout),
// 		HttpClient: &http.Client{},
// 	}

// 	type args struct {
// 		rsp   *Response
// 		input string
// 		err   error
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{
// 			name: "fine",
// 			args: args{
// 				input: "Hallo\n",
// 				rsp:   &Response{Name: "Max", Content: "Hallo"},
// 				err:   nil,
// 			},
// 		},
// 		{
// 			name: "canceled",
// 			args: args{
// 				input: "",
// 				rsp:   nil,
// 				err:   nil,
// 			},
// 		},
// 		{
// 			name: "canceled",
// 			args: args{
// 				input: "",
// 				rsp:   nil,
// 				err:   nil,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			wg := &sync.WaitGroup{}
// 			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 				jsonRsp, _ := json.Marshal(tt.args.rsp)
// 				w.Write(jsonRsp)
// 			}))

// 			ctx, cancel := context.WithCancel(context.Background())
// 			if tt.args.input == "" {
// 				c.reader = &DelayReader{delay: time.Second}
// 				go func() {
// 					time.Sleep(500 * time.Millisecond)
// 					cancel()
// 				}()
// 			}
// 			err := c.SendMessage(ts.URL, cancel, tt.args.input, wg, ctx)

// 			assert.ErrorIs(t, tt.args.err, err)

// 		})
// 	}
// }

// type DelayReader struct {
// 	delay time.Duration
// }

// func (d *DelayReader) ReadString(delim byte) (string, error) {
// 	time.Sleep(d.delay)
// 	return "you should not read this", fmt.Errorf("you should not read this")
// }

// func TestClient_ResponsePoller(t *testing.T) {
// 	var logOutput bytes.Buffer
// 	log.SetOutput(&logOutput)
// 	defer log.SetOutput(os.Stderr)

// 	c := NewClient("http://localhost:8080")
// 	c.clientId = clientId
// 	c.authToken = authToken
// 	c.Cond = sync.NewCond(c.mu)
// 	c.Endpoints = c.RegisterEndpoints(c.url)

// 	u := &UserService{
// 		chatClient: c,
// 		plugins:    RegisterPlugins(c),
// 		poll:       false,
// 		mu:         &sync.Mutex{},
// 	}

// 	u.Cond = sync.NewCond(u.mu)

// 	go u.ResponsePoller()

// 	assert.Contains(t, logOutput.String(), "")

// }

// func TestClient_extractInputToMessage(t *testing.T) {
// 	c := network.NewClient("http://localhost:8080")
// 	// anders handeln
// 	// c.clientId = clientId
// 	// c.authToken = authToken
// 	c.Endpoints = c.RegisterEndpoints(c.Url)

// 	u := NewUserService(c)

// 	type args struct {
// 		input        string
// 		outputPlugin string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{
// 			name: "/broadcast",
// 			args: args{
// 				input:        "Hi",
// 				outputPlugin: "/broadcast",
// 			},
// 		},
// 		{
// 			name: "any other plugin",
// 			args: args{
// 				input:        "/help\n",
// 				outputPlugin: "/help",
// 			},
// 		},
// 		{
// 			name: "with spaces",
// 			args: args{
// 				input:        "/users     ",
// 				outputPlugin: "/users",
// 			},
// 		},
// 		{
// 			name: "non existing",
// 			args: args{
// 				input:        "/privateABCDEFGHIJ",
// 				outputPlugin: "/privateABCDEFGHIJ",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			message := u.ParseInputToMessage(tt.args.input)
// 			assert.Equal(t, tt.args.outputPlugin, message.Plugin)
// 		})
// 	}
// }
