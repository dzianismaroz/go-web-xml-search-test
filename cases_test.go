package main

import (
	"errors"
)

type TestCase struct {
	client *SearchClient
	search SearchRequest
	result *SearchResponse
	err    error
}

var cases = []TestCase{
	//-----------------   invalid auth test cases  -------------------------
	{
		client: &SearchClient{AccessToken: ""},
		err:    errors.New("Bad AccessToken"),
	},
	{
		client: &SearchClient{AccessToken: "$GRANTMEACCESS$"},
		err:    errors.New("Bad AccessToken"),
	},
	//------------------ invalid search params --------------------------
	// {
	// 	client: &SearchClient{AccessToken: ValidToken},
	// 	search: SearchRequest{},
	// 	err:    errors.New("cant unpack result json: unexpected end of JSON input"),
	// },
	{
		client: &SearchClient{AccessToken: ValidToken},
		search: SearchRequest{Limit: -2},
		err:    errors.New("limit must be > 0"),
	},
	{
		client: &SearchClient{AccessToken: ValidToken},
		search: SearchRequest{Limit: 10, Offset: -10},
		err:    errors.New("offset must be > 0"),
	},
	{
		client: &SearchClient{AccessToken: ValidToken},
		search: SearchRequest{Limit: 10, Offset: 10, OrderBy: OrderByAsc, OrderField: "Something"},
		err:    errors.New("OrderFeld Something invalid"),
	},
	// ----------- testing timeout ---------
	// {
	// 	client: &SearchClient{AccessToken: ValidToken},
	// 	search: SearchRequest{Limit: 10, Offset: 10, OrderBy: OrderByAsc, OrderField: "Something"},
	// 	err:    errors.New("timeout for %s"),
	// },

	// ------ testing unpack json error ----------------

	// --------- success cases ---------------------
	{
		client: &SearchClient{AccessToken: ValidToken},
		search: SearchRequest{Limit: 30, OrderBy: OrderByAsc, OrderField: "Age", Query: "Boyd"},
		result: &SearchResponse{Users: []User{{Name: "Boyd Wolf", Id: 0, Age: 22,
			About: "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur " +
				"ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla " +
				"commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud " +
				"cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla" +
				" nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
			Gender: "male",
		}}},
	},

	{
		client: &SearchClient{AccessToken: ValidToken},
		search: SearchRequest{Limit: 2, OrderBy: OrderByDesc, OrderField: "Id", Query: "commodo e"},
		result: &SearchResponse{Users: []User{
			{Id: 34, Name: "Kane Sharp", Age: 34, About: "Lorem proident sint minim anim commodo cillum. Eiusmod velit culpa commodo anim consectetur consectetur sint sint labore. Mollit consequat consectetur magna nulla veniam commodo eu ut et. Ut adipisicing qui ex consectetur officia sint ut fugiat ex velit cupidatat fugiat nisi non. Dolor minim mollit aliquip veniam nostrud. Magna eu aliqua Lorem aliquip.\n", Gender: "male"},
			{Id: 1, Name: "Hilda Mayer", Age: 21, About: "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n", Gender: "female"}},
			NextPage: true},
	},

	{
		client: &SearchClient{AccessToken: ValidToken},
		search: SearchRequest{Query: ""},
		result: &SearchResponse{Users: []User{
			{Id: 0, Name: "Boyd Wolf", Age: 22, About: "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n", Gender: "male"},
			{Id: 1, Name: "Hilda Mayer", Age: 21, About: "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n", Gender: "female"},
			{Id: 2, Name: "Brooks Aguilar", Age: 25, About: "Velit ullamco est aliqua voluptate nisi do. Voluptate magna anim qui cillum aliqua sint veniam reprehenderit consectetur enim. Laborum dolore ut eiusmod ipsum ad anim est do tempor culpa ad do tempor. Nulla id aliqua dolore dolore adipisicing.\n", Gender: "male"},
			{Id: 3, Name: "Everett Dillard", Age: 27, About: "Sint eu id sint irure officia amet cillum. Amet consectetur enim mollit culpa laborum ipsum adipisicing est laboris. Adipisicing fugiat esse dolore aliquip quis laborum aliquip dolore. Pariatur do elit eu nostrud occaecat.\n", Gender: "male"},
			{Id: 4, Name: "Owen Lynn", Age: 30, About: "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n", Gender: "male"},
			{Id: 5, Name: "Beulah Stark", Age: 30, About: "Enim cillum eu cillum velit labore. In sint esse nulla occaecat voluptate pariatur aliqua aliqua non officia nulla aliqua. Fugiat nostrud irure officia minim cupidatat laborum ad incididunt dolore. Fugiat nostrud eiusmod ex ea nulla commodo. Reprehenderit sint qui anim non ad id adipisicing qui officia Lorem.\n", Gender: "female"},
			{Id: 6, Name: "Jennings Mays", Age: 39, About: "Veniam consectetur non non aliquip exercitation quis qui. Aliquip duis ut ad commodo consequat ipsum cupidatat id anim voluptate deserunt enim laboris. Sunt nostrud voluptate do est tempor esse anim pariatur. Ea do amet Lorem in mollit ipsum irure Lorem exercitation. Exercitation deserunt adipisicing nulla aute ex amet sint tempor incididunt magna. Quis et consectetur dolor nulla reprehenderit culpa laboris voluptate ut mollit. Qui ipsum nisi ullamco sit exercitation nisi magna fugiat anim consectetur officia.\n", Gender: "male"},
			{Id: 7, Name: "Leann Travis", Age: 34, About: "Lorem magna dolore et velit ut officia. Cupidatat deserunt elit mollit amet nulla voluptate sit. Quis aute aliquip officia deserunt sint sint nisi. Laboris sit et ea dolore consequat laboris non. Consequat do enim excepteur qui mollit consectetur eiusmod laborum ut duis mollit dolor est. Excepteur amet duis enim laborum aliqua nulla ea minim.\n", Gender: "female"},
			{Id: 8, Name: "Glenn Jordan", Age: 29, About: "Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n", Gender: "male"},
			{Id: 9, Name: "Rose Carney", Age: 36, About: "Voluptate ipsum ad consequat elit ipsum tempor irure consectetur amet. Et veniam sunt in sunt ipsum non elit ullamco est est eu. Exercitation ipsum do deserunt do eu adipisicing id deserunt duis nulla ullamco eu. Ad duis voluptate amet quis commodo nostrud occaecat minim occaecat commodo. Irure sint incididunt est cupidatat laborum in duis enim nulla duis ut in ut. Cupidatat ex incididunt do ullamco do laboris eiusmod quis nostrud excepteur quis ea.\n", Gender: "female"},
			{Id: 10, Name: "Henderson Maxwell", Age: 30, About: "Ex et excepteur anim in eiusmod. Cupidatat sunt aliquip exercitation velit minim aliqua ad ipsum cillum dolor do sit dolore cillum. Exercitation eu in ex qui voluptate fugiat amet.\n", Gender: "male"},
			{Id: 11, Name: "Gilmore Guerra", Age: 32, About: "Labore consectetur do sit et mollit non incididunt. Amet aute voluptate enim et sit Lorem elit. Fugiat proident ullamco ullamco sint pariatur deserunt eu nulla consectetur culpa eiusmod. Veniam irure et deserunt consectetur incididunt ad ipsum sint. Consectetur voluptate adipisicing aute fugiat aliquip culpa qui nisi ut ex esse ex. Sint et anim aliqua pariatur.\n", Gender: "male"},
			{Id: 12, Name: "Cruz Guerrero", Age: 36, About: "Sunt enim ad fugiat minim id esse proident laborum magna magna. Velit anim aliqua nulla laborum consequat veniam reprehenderit enim fugiat ipsum mollit nisi. Nisi do reprehenderit aute sint sit culpa id Lorem proident id tempor. Irure ut ipsum sit non quis aliqua in voluptate magna. Ipsum non aliquip quis incididunt incididunt aute sint. Minim dolor in mollit aute duis consectetur.\n", Gender: "male"},
			{Id: 13, Name: "Whitley Davidson", Age: 40, About: "Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n", Gender: "male"},
			{Id: 14, Name: "Nicholson Newman", Age: 23, About: "Tempor minim reprehenderit dolore et ad. Irure id fugiat incididunt do amet veniam ex consequat. Quis ad ipsum excepteur eiusmod mollit nulla amet velit quis duis ut irure.\n", Gender: "male"},
			{Id: 15, Name: "Allison Valdez", Age: 21, About: "Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n", Gender: "male"},
			{Id: 16, Name: "Annie Osborn", Age: 35, About: "Consequat fugiat veniam commodo nisi nostrud culpa pariatur. Aliquip velit adipisicing dolor et nostrud. Eu nostrud officia velit eiusmod ullamco duis eiusmod ad non do quis.\n", Gender: "female"},
			{Id: 17, Name: "Dillard Mccoy", Age: 36, About: "Laborum voluptate sit ipsum tempor dolore. Adipisicing reprehenderit minim aliqua est. Consectetur enim deserunt incididunt elit non consectetur nisi esse ut dolore officia do ipsum.\n", Gender: "male"},
			{Id: 18, Name: "Terrell Hall", Age: 27, About: "Ut nostrud est est elit incididunt consequat sunt ut aliqua sunt sunt. Quis consectetur amet occaecat nostrud duis. Fugiat in irure consequat laborum ipsum tempor non deserunt laboris id ullamco cupidatat sit. Officia cupidatat aliqua veniam et ipsum labore eu do aliquip elit cillum. Labore culpa exercitation sint sint.\n", Gender: "male"},
			{Id: 19, Name: "Bell Bauer", Age: 26, About: "Nulla voluptate nostrud nostrud do ut tempor et quis non aliqua cillum in duis. Sit ipsum sit ut non proident exercitation. Quis consequat laboris deserunt adipisicing eiusmod non cillum magna.\n", Gender: "male"},
			{Id: 20, Name: "Lowery York", Age: 27, About: "Dolor enim sit id dolore enim sint nostrud deserunt. Occaecat minim enim veniam proident mollit Lorem irure ex. Adipisicing pariatur adipisicing aliqua amet proident velit. Magna commodo culpa sit id.\n", Gender: "male"},
			{Id: 21, Name: "Johns Whitney", Age: 26, About: "Elit sunt exercitation incididunt est ea quis do ad magna. Commodo laboris nisi aliqua eu incididunt eu irure. Labore ullamco quis deserunt non cupidatat sint aute in incididunt deserunt elit velit. Duis est mollit veniam aliquip. Nulla sunt veniam anim et sint dolore.\n", Gender: "male"},
			{Id: 22, Name: "Beth Wynn", Age: 31, About: "Proident non nisi dolore id non. Aliquip ex anim cupidatat dolore amet veniam tempor non adipisicing. Aliqua adipisicing eu esse quis reprehenderit est irure cillum duis dolor ex. Laborum do aute commodo amet. Fugiat aute in excepteur ut aliqua sint fugiat do nostrud voluptate duis do deserunt. Elit esse ipsum duis ipsum.\n", Gender: "female"},
			{Id: 23, Name: "Gates Spencer", Age: 21, About: "Dolore magna magna commodo irure. Proident culpa nisi veniam excepteur sunt qui et laborum tempor. Qui proident Lorem commodo dolore ipsum.\n", Gender: "male"},
			{Id: 24, Name: "Gonzalez Anderson", Age: 33, About: "Quis consequat incididunt in ex deserunt minim aliqua ea duis. Culpa nisi excepteur sint est fugiat cupidatat nulla magna do id dolore laboris. Aute cillum eiusmod do amet dolore labore commodo do pariatur sit id. Do irure eiusmod reprehenderit non in duis sunt ex. Labore commodo labore pariatur ex minim qui sit elit.\n", Gender: "male"},
			{Id: 25, Name: "Katheryn Jacobs", Age: 32, About: "Magna excepteur anim amet id consequat tempor dolor sunt id enim ipsum ea est ex. In do ea sint qui in minim mollit anim est et minim dolore velit laborum. Officia commodo duis ut proident laboris fugiat commodo do ex duis consequat exercitation. Ad et excepteur ex ea exercitation id fugiat exercitation amet proident adipisicing laboris id deserunt. Commodo proident laborum elit ex aliqua labore culpa ullamco occaecat voluptate voluptate laboris deserunt magna.\n", Gender: "female"},
			{Id: 26, Name: "Sims Cotton", Age: 39, About: "Ex cupidatat est velit consequat ad. Tempor non cillum labore non voluptate. Et proident culpa labore deserunt ut aliquip commodo laborum nostrud. Anim minim occaecat est est minim.\n", Gender: "male"},
			{Id: 27, Name: "Rebekah Sutton", Age: 26, About: "Aliqua exercitation ad nostrud et exercitation amet quis cupidatat esse nostrud proident. Ullamco voluptate ex minim consectetur ea cupidatat in mollit reprehenderit voluptate labore sint laboris. Minim cillum et incididunt pariatur amet do esse. Amet irure elit deserunt quis culpa ut deserunt minim proident cupidatat nisi consequat ipsum.\n", Gender: "female"},
			{Id: 28, Name: "Cohen Hines", Age: 32, About: "Deserunt deserunt dolor ex pariatur dolore sunt labore minim deserunt. Tempor non et officia sint culpa quis consectetur pariatur elit sunt. Anim consequat velit exercitation eiusmod aute elit minim velit. Excepteur nulla excepteur duis eiusmod anim reprehenderit officia est ea aliqua nisi deserunt officia eiusmod. Officia enim adipisicing mollit et enim quis magna ea. Officia velit deserunt minim qui. Commodo culpa pariatur eu aliquip voluptate culpa ullamco sit minim laboris fugiat sit.\n", Gender: "male"},
			{Id: 29, Name: "Clarissa Henry", Age: 34, About: "Nostrud enim ea ad reprehenderit tempor ullamco exercitation. Elit in voluptate pariatur sit nisi occaecat laboris esse ipsum. Mollit elit et deserunt ea laboris sunt est amet culpa laboris occaecat ipsum sunt sunt.\n", Gender: "female"},
			{Id: 30, Name: "Dickson Silva", Age: 32, About: "Ipsum aliqua proident ullamco laboris eu occaecat deserunt. Amet ut adipisicing sint veniam dolore aliquip est mollit ex officia esse eiusmod veniam. Dolore magna minim aliquip sit deserunt. Nostrud occaecat dolore aliqua aliquip voluptate aliquip ad adipisicing.\n", Gender: "male"},
			{Id: 31, Name: "Palmer Scott", Age: 37, About: "Elit fugiat commodo laborum quis eu consequat. In velit magna sit fugiat non proident ipsum tempor eu. Consectetur exercitation labore eiusmod occaecat adipisicing irure consequat fugiat ullamco aliquip nostrud anim irure enim. Duis do amet cillum eiusmod eu sunt. Minim minim sunt sit sit enim velit sint tempor enim sint aliquip voluptate reprehenderit officia. Voluptate magna sit consequat adipisicing ut eu qui.\n", Gender: "male"},
			{Id: 32, Name: "Christy Knapp", Age: 40, About: "Incididunt culpa dolore laborum cupidatat consequat. Aliquip cupidatat pariatur sit consectetur laboris labore anim labore. Est sint ut ipsum dolor ipsum nisi tempor in tempor aliqua. Aliquip labore cillum est consequat anim officia non reprehenderit ex duis elit. Amet aliqua eu ad velit incididunt ad ut magna. Culpa dolore qui anim consequat commodo aute.\n", Gender: "female"},
			{Id: 33, Name: "Twila Snow", Age: 36, About: "Sint non sunt adipisicing sit laborum cillum magna nisi exercitation. Dolore officia esse dolore officia ea adipisicing amet ea nostrud elit cupidatat laboris. Proident culpa ullamco aute incididunt aute. Laboris et nulla incididunt consequat pariatur enim dolor incididunt adipisicing enim fugiat tempor ullamco. Amet est ullamco officia consectetur cupidatat non sunt laborum nisi in ex. Quis labore quis ipsum est nisi ex officia reprehenderit ad adipisicing fugiat. Labore fugiat ea dolore exercitation sint duis aliqua.\n", Gender: "female"},
			{Id: 34, Name: "Kane Sharp", Age: 34, About: "Lorem proident sint minim anim commodo cillum. Eiusmod velit culpa commodo anim consectetur consectetur sint sint labore. Mollit consequat consectetur magna nulla veniam commodo eu ut et. Ut adipisicing qui ex consectetur officia sint ut fugiat ex velit cupidatat fugiat nisi non. Dolor minim mollit aliquip veniam nostrud. Magna eu aliqua Lorem aliquip.\n", Gender: "male"}},
			NextPage: false,
		},
	},
}
