# ltvpredictor
## BY
Рэпазіторый зьмяшчае 2 версіі тэставага заданьня. Абедзьве версіі працуюць аднолькава, аднак выбар якую карыстаць залежыць ад таго, што патрэбна карыстальніку/кліенту/Мікалаю: прылада, якую легка падтрымліваць і распрацоўваць далей, ці проста кансольная прылада, якая дае адзін вынік. Мне падабаецца стуктура і ўзаемаадносіны першага варыянта, але, як я і казаў, ўсе залежыць ад патрабаваньняў (якія не былі дастаткова выразнымі на гэты сэнс, таму я вырашыў зрабіць дзьве версіі)

- letpredictor_base - першая версія, якая зьмяшчае пашыраны функцыянальнасць:
  - паўнацэнныя многія-да-многіх адносіны, існасьці Кампаніі і Краіны, якія зьмяшчаюць інфармацыю пра Аналітыку (спасылкі на існасьці, каб не марнаваць пямяць)
  - маючы на ўвазе што напісана вышей, мы можам пашыраць функцыянальнасьць як захочам, рабіць кросс спасылкі, глядзець як кампаніі паводзяць сябе ў краінах і наадварот, і гэтак далей
  - есьць магчымасьць будаваць графікі, з дадзеных, якія мы атрымалі. Прадказанне робіцца на кожны дзень да 60га, для візуалізацыі дадзеных карыстаецца бібліятэка go-echart, графікі захоўваюцца ў папку

- ltvpredictor_light - палегчаная версія ltvpredictor_base
  - няма суадносін многія-да-многіх, прылада толькі зьбірае дадзеныя для кампаній/краін ў залежнасьці ад параметра aggregate які быў перададзены
  - краіна і кампанія цяпер гэта адна існасьць, якая мае назву, сумірванныя дадзеныя аналітыкі і колькасьць карыстальнікаў
  - гэта дазваляе сэканоміць  памяць пакуль мы распазнаем файлы, але ж код ад гэтага становіцца менш чытабельным і зразумелым

Абедзьве прылады працуюць па-аднолькаваму алгарытму: правяраюць ўвод, ініцыялізуюць парсер у залежнасьці ад тыпу файла, распазнаюць файл у нашыя існасьці, і робім прадказанне па атрыманых дадзеных. 
Кожнае прадказаньне робіцца ў сваёй руціне, рэзультаты прабрасваюцца наверх ў галоўныю клясу, дзе яны выводзяцца ў кансоль.

Інфармацыю па ўсім наладкам прылады можна атрымаць праз запуск: **go run . -h**

Запуск **go run .** стартуе прыладу з параметрамі па ўмаўчанню

Поўны фармат каманды для прылады: **go run . -aggregate country -source test_data.json -model quad -graph 1**


## EN
The repo contains the 2 versions of the requested utility. They both are valid and produce same results, however it all comes down to what is needed by the end user/customer/Mikalai: the utility that will be used and further developed, or just a straignt on point one result generator. I personally like the structure of the first one, however in the end it all comes down to the requirements (that were vague in that regard, hence I decided that 2 versions should be made)

- ltvpredictor_base - the first version, containing extended functionality such as:
  - full many to many relationship management, with Country and Campaigns entities holding information about Analytics (links to objects, so we don't double the data)
  - with mentioned above we can extend functionality where needed, making cross references, seeing how certain campaigns behave in different countries, and such
  - ability to build graphs on predicted data. I predict all the days up to day 60, and use go-echart to visualize the data. charts are stored in the separate folder

- ltvpredictor_light - the stripped down version of the ltvpredictor_base
  - no many to many relations, the utility is on point and just gathers the data that we need, depending on the values of the aggregate flag passed to the untility
  - country and campaign are now single entity, having name, summed up analytics and usercount
  - this allows to save memory while parsing the data, however makes the code less readable and understandable

Both applications function in a similar fashion: check the input for validity, init the parser basing on detected file type, parse the files into corresponding structures, perform prediction on the parsed data.
Each prediction is done in a separate routine, results are passed back to main app where they are printed to the console.

All the information about how the utility functions can be obtained with: **go run . -h**

Running with **go run .** starts the app with a default values.

Full format of the command to run the utility:  **go run . -aggregate country -source test_data.json -model quad -graph 1**
