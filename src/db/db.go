package config

import (
	"database/sql"
)

var DB *sql.DB

const SqlGEN3 = `select count(id) as prod, COALESCE(MAX(LEFT(RTRIM(CodeCeHousing), 1)),'')  as model
					from [GEN3_Clone].[dbo].[GEN3_LINE_B_PROD]
					WHERE LineB2Good=1 and LineNumber=2 and DateDataSave >= @dataInicial and DateDataSave <= @dataFinal`

const SqlInv3 = `select COUNT(MINDEX) as prod, COALESCE(MAX(LEFT(RIGHT(RTRIM(PART_NUMBER), 3),1)),'') as model
					 FROM [INVERTER3_Clone].[dbo].[INVERTER_PROD]
					 where REJCODE is null and WRITE_STATION = '140' and TIME_STAMP >= @dataInicial and TIME_STAMP <= @dataFinal`

const SqlR744 = `select COUNT(MINDEX) as prod, COALESCE(MAX(RIGHT(LEFT(RTRIM(PN_END_ITEM), 9),4)),'') as model
					 FROM [R744_Clone].[dbo].[LineD]
					 where (REJECTED IS NULL OR REJECTED = 0) and WRITE_STATION = '670' and TIME_STAMP >= @dataInicial and TIME_STAMP <= @dataFinal`

const SqlYF = `select COUNT(MINDEX) as prod, COALESCE(MAX(RIGHT(LEFT(RTRIM(END_ITEM_PART_NUMBER), 10),5)),'') as model
				 FROM [YF_Clone].[dbo].[LineD]
				 where REJECTED = 0 and WRITE_STATION = 990 and TIME_STAMP >= @dataInicial and TIME_STAMP <= @dataFinal`

const SqlInv4 = `select COUNT(MINDEX) as prod, COALESCE(MAX(RIGHT(LEFT(RTRIM(PN_INVERTER), 8),3)),'') as model
				 FROM [GEN4_Inverter].[dbo].[INVERTER_PROD]
				 where (REJECTED IS NULL OR REJECTED = 0) and WRITE_STATION = '90' and TIME_STAMP >= @dataInicial and TIME_STAMP <= @dataFinal`

const SqlInv42 = `select COUNT(MINDEX) as prod, COALESCE(MAX(RIGHT(LEFT(RTRIM(PN_INVERTER), 8),3)),'') as model
				 FROM [GEN4_Inverter].[dbo].[INVERTER_LINE2_PROD]
				 where (REJECTED IS NULL OR REJECTED = 0) and WRITE_STATION = '90' and TIME_STAMP >= @dataInicial and TIME_STAMP <= @dataFinal`

const SqlInv43 = `select COUNT(MINDEX) as prod, COALESCE(MAX(RIGHT(LEFT(RTRIM(PN_INVERTER), 8),3)),'') as model
				 FROM [GEN4_Inverter].[dbo].[INVERTER_LINE3_PROD]
				 where (REJECTED IS NULL OR REJECTED = 0) and WRITE_STATION = '90' and TIME_STAMP >= @dataInicial and TIME_STAMP <= @dataFinal`

const SqlmodelCheck = `SELECT COALESCE([name],'')
                     FROM [TESTEProd].[dbo].[models] 
                     where CHARINDEX(@model, fassy_models) > 0 and line_id = @line_id`

const SqlInsertTicket = `INSERT INTO [TESTEProd].dbo.[ScrapTickets] (ticket,date,lastupdated,price,costCenter) VALUES (@ticket,@date,@lastupdated,@price,@costcenter)`

const SqlInsertTicketNums = `INSERT INTO [TESTEProd].dbo.[ScrapTickets] (ticket,date,requester,person,lastupdated,price) VALUES (?,?,?,?,?,?)`

const SqlUpdatePerson = `UPDATE [TESTEProd].dbo.[ScrapTickets] SET person = @person,lastupdated = @lastupdated WHERE ticket = @ticket`

const SqlUpdateRequester = `UPDATE [TESTEProd].dbo.[ScrapTickets] SET requester = @person,lastupdated = @lastupdated WHERE ticket = @ticket`

const SqlHeartbeat = `INSERT INTO [TESTEProd].dbo.[com_heartbeat] (machine,ip,app,timestamp) VALUES (@machine,@ip,@app,@timestamp)`
