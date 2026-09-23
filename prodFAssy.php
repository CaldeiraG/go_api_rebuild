<?php

$stmt = $conn->query("SELECT [Name] FROM [TESTEProd].[dbo].[lines] where id=$line");
                                    
while ($row = $stmt->fetch())
{
    $name = $row['Name'];
}

if($model != "")
{

	if($model == "all")
	{
		$model = "";
	}
	else
	{

        $modelFAssy = $model;

		/*$stmt2 = $conn->query("SELECT [fassy_models] FROM [TESTEProd].[dbo].[models] where fassy_models like $model");
											
		while ($row = $stmt2->fetch())
		{
			$modelFAssy = $row['fassy_models'];

			//echo $modelFAssy;
		}*/
	}
}

//$paramModel = ($model == "" ? $paramModel = "" : $paramModel = "");

switch ($line)
		{
					// *** VS14 databases *** \\

		case "42":  //$databaseInUse = "[VS14_Clone].[dbo].[VS14_LINE_A_PROD]"; 
					//$name = "VS14 Line A";
					$databaseInUse = "[VS14LineA_Production].[dbo].[VS14AResults]"; 
					$ID = "ID"; 
					$dateTime = "DateTime";
					$param = "GlobalResult=1";
					$paramRej = "GlobalResult=0";
					$paramRejSta = "GlobalResult=0 AND rejected_station='$station'";
					$paramModel = ($model != "" ? $paramModel = "AND model_id in ($modelFAssy) AND ": $paramModel = "AND");
					$model_id  = "";
					
					break;
			
		case "43":  //$databaseInUse = "[VS14_Clone].[dbo].[VS14_LINE_B_PROD]";
					//$name = "VS14 Line B"; 
					$databaseInUse = "[VS14LineB_Production].[dbo].[VS14BResultValues]";
					$ID = "ID";
					$dateTime = "DateTime";
					$param = "GlobalResult=1";
					$paramRej = "GlobalResult=0";
					$paramRejSta = "GlobalResult=0 AND rejected_station='$station'";
					$paramModel = ($model != "" ? $paramModel = "AND model_id in ($modelFAssy) AND ": $paramModel = "AND");
					$model_id  = "";

					break;
		   
		case "44":  //$databaseInUse = "[VS14_Clone].[dbo].[VS14_LINE_C_PROD]";
					//$name = "VS14 Line C";
					$databaseInUse = "[FM12240_Visteon].[dbo].[FinishedPallets]";
					$ID = "ID"; 
					$dateTime = "DateTimeInsertedOnDB";
					$param = "Rejected=0";
					$paramRej = "Rejected=1";
					$paramRejSta = "Rejected=1 AND RejStation=".(is_numeric($station) ? ($station / 10) : $station)."";
					$paramModel = ($model != "" ? $paramModel = "AND model_id in ($modelFAssy) AND ": $paramModel = "AND");
					$model_id  = "";
				break;

					// *** GEN 3 databases *** \\

		case "45":  include 'db_connection_graph.php'; //$databaseInUse = "[GEN3_Clone].[dbo].[GEN3_LINE_A_PROD]"; 
					$databaseInUse = "[GEN3_Clone].[dbo].[GEN3_LINE_A_PROD]";
					//$name = "Gen3.8 Line A",;
					$ID = "ID"; 
					$dateTime = "Data";
					$param = "PALETE_DB_STATUS_OK_NOK like '1'";
					$paramRej = "PALETE_DB_STATUS_OK_NOK like '0'";
					$paramRejSta = "PALETE_DB_STATUS_OK_NOK like '0' AND PALETE_DB_STATUS_STF like '".$station."'";
					$paramModel = ($model != "" ? $paramModel = "AND PALETE_DB_SN_SNCH like '$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "PALETE_DB_SN_SNCH";
				break;

		case "46":  include 'db_connection_graph.php'; //$databaseInUse = "[GEN3_Clone].[dbo].[GEN3_LINE_B_PROD]";
					//$name = "Gen3.8 Line B1";
					$databaseInUse = "[GEN3_Clone].[dbo].[GEN3_LINE_B_PROD]";
					$ID = "ID"; 
					$dateTime = "DateDataSave";
					$param = "LineNumber=1 AND LineB1Good=1";
					$paramRej = "LineNumber=1 AND LineB1Good=0";
					$paramRejSta = "LineNumber=1 AND LineB1Good=0 AND RejStation=".$station."";
					$paramModel = ($model != "" ? $paramModel = "AND CodeCeHousing like '$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "CodeCeHousing";
				break;

		case "47":  include 'db_connection_graph.php'; //$databaseInUse = "[GEN3_Clone].[dbo].[GEN3_LINE_B_PROD]";
					//$name = "Gen3.8 Line B2";
					$databaseInUse = "[GEN3_Clone].[dbo].[GEN3_LINE_B_PROD]";
					$ID = "ID"; 
					$dateTime = "DateDataSave";
					$param = "LineNumber=2 AND LineB2Good=1";
					$paramRej = "LineNumber=2 AND LineB2Good=0";
					$paramRejSta = "LineNumber=2 AND LineB2Good=0 AND RejStation=".$station."";
					$paramModel = ($model != "" ? $paramModel = "AND CodeCeHousing like '$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "CodeCeHousing";
				break;

					// *** YF databases *** \\

		case "49":  include 'db_connection_yflinea.php';
		
					$databaseInUse = "[NewHanon].[dbo].[hanon_result_table]";
					//$name = "YF Line A";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "WRITE_STATION = 290 AND REJ_STATION = 0";
					$paramRej = "REJ_STATION > 0";
					$paramRejSta = "REJECTED = 1 AND REJ_STATION like '".$station."'";
					$paramModel = ($model != "" ? $paramModel = "AND PART_NUMBER like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "PART_NUMBER";
					;

					break;

		case "50":  include 'db_connection_yflineb.php';
		
					$databaseInUse = "[NewHanon].[dbo].[hanon_result_table]";
					//$name = "YF Line B";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "REJECTED = 0 and WRITE_STATION = 460";
					$paramRej = "REJECTED = 1";
					$paramRejSta = "REJECTED = 1 AND REJ_STATION like '".$station."'";
					$paramModel = ($model != "" ? $paramModel = "AND PART_NUMBER like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "PART_NUMBER";
					;

					break;

		case "51":  include 'db_connection_yflinec.php';
		
					$databaseInUse = "[NewHanon].[dbo].[hanon_result_table]";
					//$name = "YF Line C";
					$ID = "MINDEX";
					$dateTime = "TIME_STAMP";
					$param = "REJECTED = 0 and WRITE_STATION = 670";
					$paramRej = "REJECTED > 0";
					$paramRejSta = "REJECTED > 0 AND REJSTATION like '".$station."'";
					$paramModel = ($model != "" ? $paramModel = "AND PART_NUMBER like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "PART_NUMBER";
					;

					break;

		case "52":  include 'db_connection_yflined.php';
		
					$databaseInUse = "[NewHanon].[dbo].[hanon_result_table]";
					//$name = "YF Line D";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "REJECTED = 0 and WRITE_STATION = 990";
					$paramRej = "REJECTED = 1";
					$paramRejSta = "REJECTED = 1 AND REJECT_STATION like '".$station."'";
					$paramModel = ($model != "" ? $paramModel = "AND END_ITEM_PART_NUMBER like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "END_ITEM_PART_NUMBER";
					;

					break;

					// *** Gen3.8 Inverter *** \\

		case "53":  include 'db_connection_sql01.php';
		
					$databaseInUse = "[Inverter_Clone].[dbo].[Line1]";
					//$name = "Gen3.8 Inverter";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "AND REJCODE IS NULL and WRITE_STATION = '140'";
					$paramRej = "AND REJCODE > 0 AND REJSTATION !=90";
					$paramRejSta = "AND REJCODE > 0 AND REJSTATION like '".$station."'";

                    if($model != "null" && $model != "BMW" && $model != "VW")
                    {
                        $paramModel = "AND RIGHT(RTRIM(PART_NUMBER), 4) like '%$modelFAssy%'";
                    }
                    else if ($model == "BMW")
                    {
                        $paramModel = "AND ((RIGHT(RTRIM(PART_NUMBER), 4)) like '%T%'
                                        OR (RIGHT(RTRIM(PART_NUMBER), 4)) like '%U%'
                                        OR (RIGHT(RTRIM(PART_NUMBER), 4)) like '%V%'
                                        OR (RIGHT(RTRIM(PART_NUMBER), 4)) like '%Y%'
                                        OR (RIGHT(RTRIM(PART_NUMBER), 4)) like '%Z%')";
                    }
                    else if ($model == "VW")
                    {
                        $paramModel = "AND ((RIGHT(RTRIM(PART_NUMBER), 4)) like '%A%'
                                        OR (RIGHT(RTRIM(PART_NUMBER), 4)) like '%B%'
                                        OR (RIGHT(RTRIM(PART_NUMBER), 4)) like '%C%')";
                    }
                    $model_id  = "model_id";
				break;

					// *** R744 databases *** \\

		case "54":  include 'db_connection_r744linea.php';
		
					$databaseInUse = "[R744ALINE].[dbo].[hanon_result_table]";
					//$name = "R744 Line A";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "WRITE_STATION = '250' AND (REJECTED IS NULL OR REJECTED = 0)";
					$paramRej = "REJECTED = 1";
					$paramRejSta = "REJECTED = 1 AND REJECT_STATION like '".$station."'";
					$paramModel = ($model != "" ? $paramModel = "AND model_id in ($modelFAssy) AND ": $paramModel = "AND");
					$model_id  = "model_id";
				break;
		
		case "55":  include 'db_connection_r744linea.php';
		
					$databaseInUse = "[R744ALINE].[dbo].[hanon_motorhousing_table]";
					//$name = "R744 Line A";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "WRITE_STATION = '72' AND (REJECTED IS NULL OR REJECTED = 0)";
					$paramRej = "REJECTED = 1";
					$paramRejSta = "REJECTED = 1 AND REJECT_STATION like '".$station."'";
					$paramModel = ($model != "" ? $paramModel = "AND model_id in ($modelFAssy) AND ": $paramModel = "AND");
					$model_id  = "model_id";
				break;
		
		
		case "56":  include 'db_connection_r744lineb.php';
		
					$databaseInUse = "[R744BLINE].[dbo].[hanon_result_table]";
					//$name = "R744 Line B";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "WRITE_STATION = '390' AND (REJECTED IS NULL OR REJECTED = 0)";
					$paramRej = "REJECT_STATION > 0";
					$paramRejSta = "REJECT_STATION > 0 and REJECT_STATION like '".$station."'";
					$paramModel = ($model != "" ? $paramModel = "AND model_id in ($modelFAssy) AND ": $paramModel = "AND");
					$model_id  = "model_id";
				break;
		case "57":  include 'db_connection_graph.php';
		
					$databaseInUse = "[R744_Clone].[dbo].[LineC]";
					//$name = "R744 Line C";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "WRITE_STATION = '510' AND (REJECTED IS NULL OR REJECTED = 0)";
					$paramRej = "REJECTED = 1";
					$paramRejSta = "REJECTED = 1 AND REJSTATION like '".$station."'";
					$paramModel = ($model != "" ? $paramModel = "AND PART_NUMBER like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "";
				break;
		case "90":  include 'db_connection_graph.php';
		
					$databaseInUse = "[R744_Clone].[dbo].[LineD]";
					//$name = "R744 Line D";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "WRITE_STATION = '670' AND (REJECTED IS NULL OR REJECTED = 0)";
					$paramRej = "REJECTED = 1";
					$paramRejSta = "REJECTED = 1 AND REJECT_STATION like '".$station."'";
					$paramModel = ($model != "" ? $paramModel = "AND PN_END_ITEM like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "model_id";

					break;
					
		case "83":  include 'db_connection_graph.php';
		
					$databaseInUse = "[GEN4_Inverter].[dbo].[INVERTER_PROD]";
					//$name = "R744 Line D";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "WRITE_STATION = '90' AND (REJECTED IS NULL OR REJECTED = 0)";
					$paramRej = "REJECTED = 1";
					$paramRejSta = "REJECTED = 1 AND REJECT_STATION like '".$station."'";
					$paramModel = ($model != "" ? $paramModel = "AND PN_INVERTER like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "PN_INVERTER";
									break;
					
		case "84":  //include 'db_connection_graph.php';
		
					$error = "";
					exit($error);
					break;
					
		case "85":  //include 'db_connection_graph.php';
		
					include 'db_connection_graph.php';
		
					$databaseInUse = "[GEN4_Inverter].[dbo].[INVERTER_LINE2_PROD]";
					//$name = "R744 Line D";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "WRITE_STATION = '90' AND (REJECTED IS NULL OR REJECTED = 0)";
					$paramRej = "REJECTED = 1";
					$paramRejSta = "REJECTED = 1 AND REJECT_STATION like '".$station."'";
					$paramModel = ($model != "" ? $paramModel = "AND PN_INVERTER like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "PN_INVERTER";
									break;
		case "86":  //include 'db_connection_graph.php';
		
					$error = "";
					exit($error);
					break;
		case "87":  //include 'db_connection_graph.php';
		
					$error = "";
					exit($error);
					break;
		case "88":  //include 'db_connection_graph.php';
		
					$error = "";
					exit($error);
					break;
		case "89":  //include 'db_connection_graph.php';
		
					$error = "";
					exit($error);
					break;

		case "91":  //include 'db_connection_graph.php';
		
					include 'db_connection_graph.php';
		
					$databaseInUse = "[GEN4_Inverter].[dbo].[INVERTER_LINE3_PROD]";
					//$name = "R744 Line D";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "WRITE_STATION = '90' AND (REJECTED IS NULL OR REJECTED = 0)";
					$paramRej = "REJECTED = 1";
					$paramRejSta = "REJECTED = 1 AND REJECT_STATION like '".$station."'";
					$paramModel = ($model != "" ? $paramModel = "AND PN_INVERTER like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "PN_INVERTER";
									break;
									
		case "92":  //include 'db_connection_graph.php';
		
					include 'db_connection_graph.php';
		
					$databaseInUse = "[R744_Clone].[dbo].[BPV_ProductionData]";
					//$name = "R744 Line D";
					$ID = "ID"; 
					$dateTime = "DateCreated";
					$param = "RESULT = 1";
					$paramRej = "RESULT = 2";
					$paramRejSta = "RESULT = 2";
					$paramModel = ($model != "" ? $paramModel = "AND RIGHT(RTRIM(PN_INVERTER), 7) like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "model_id";
									break;
					// *** Não existe *** \\
					
		case "96":  //include 'db_connection_graph.php';
		
					include 'db_connection_graph.php';
		
					$databaseInUse = "[GEN4_Inverter].[dbo].[INVERTER_LINE3_SOLDERING]";
					//$name = "R744 Line D";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "WRITE_STATION = '9'";
					$paramRej = "CYCLE_TIME = 0";
					$paramRejSta = "";
					$paramModel = ($model != "" ? $paramModel = "AND PN_INVERTER like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "BARCODE_1";
									break;
									
								
		case "97":  //include 'db_connection_graph.php';
		
					include 'db_connection_graph.php';
		
					$databaseInUse = "[GEN4_Inverter].[dbo].[INVERTER_LINE1_SOLDERING]";
					//$name = "R744 Line D";
					$ID = "MINDEX"; 
					$dateTime = "TIME_STAMP";
					$param = "WRITE_STATION = '9'";
					$paramRej = "ST09_CYCLE_TIME = 0";
					$paramRejSta = "";
					$paramModel = ($model != "" ? $paramModel = "AND PN_INVERTER like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "ST09_BARCODE_1";
									break;
									
		case "1096":  //include 'db_connection_graph.php';
		
					include 'db_connection_graph.php';
		
					$databaseInUse = "[YF_Clone].[dbo].[BRCKT_BMW_Datetime]";
					//$name = "R744 Line D";
					$ID = "ID"; 
					$dateTime = "Date";
					$param = "ErrorCode != 0";
					$paramRej = "ErrorCode = 0";
					$paramRejSta = "";
					$paramModel = ($model != "" ? $paramModel = "AND PN_INVERTER like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "SN_Compressor";
									break;
									
		case "1097":  //include 'db_connection_graph.php';
		
					include 'db_connection_graph.php';
		
					$databaseInUse = "[YF_Clone].[dbo].[BRCKT_VW]";
					//$name = "R744 Line D";
					$ID = "ID"; 
					$dateTime = "DateTime";
					$param = "Result = 1 and TaskName = 'SCREW 1'";
					$paramRej = "Result = 0 and TaskName = 'SCREW 1'";
					$paramRejSta = "";
					$paramModel = ($model != "" ? $paramModel = "AND PN_INVERTER like '%$modelFAssy%' AND ": $paramModel = "AND");
					$model_id  = "partNum";
									break;
        case "1098":  include 'db_connection_graph.php';

                    $databaseInUse = "[YF_Clone].[dbo].[hanon_chrotor_table]";
                    //$name = "YF Line A";
                    $ID = "MINDEX";
                    $dateTime = "TIME_STAMP";
                    $param = "REJ_STATION = 0";
                    $paramRej = "REJECTED = 1";
                    $paramRejSta = "REJECTED = 1 AND REJ_STATION like '".$station."'";
                    $paramModel = ($model != "" ? $paramModel = "AND PART_NUMBER like '%$modelFAssy%' AND ": $paramModel = "AND");
                    $model_id  = "CH_BARCODE";
                    ;

                    break;
        case "1099":  include 'db_connection_graph.php';

                    $databaseInUse = "[YF_Clone].[dbo].[hanon_fiexed_table]";
                    //$name = "YF Line A";
                    $ID = "MINDEX";
                    $dateTime = "TIME_STAMP";
                    $param = "REJ_STATION = 0";
                    $paramRej = "REJECTED = 1";
                    $paramRejSta = "REJECTED = 1 AND REJ_STATION like '".$station."'";
                    $paramModel = ($model != "" ? $paramModel = "AND PART_NUMBER like '%$modelFAssy%' AND ": $paramModel = "AND");
                    $model_id  = "FS_BARCODE";
                    ;

                    break;
        case "1100":  include 'db_connection_graph.php';

                    $databaseInUse = "[YF_Clone].[dbo].[hanon_orbiting_table]";
                    //$name = "YF Line A";
                    $ID = "MINDEX";
                    $dateTime = "TIME_STAMP";
                    $param = "REJ_STATION = 0";
                    $paramRej = "REJECTED = 1";
                    $paramRejSta = "REJECTED = 1 AND REJ_STATION like '".$station."'";
                    $paramModel = ($model != "" ? $paramModel = "AND PART_NUMBER like '%$modelFAssy%' AND ": $paramModel = "AND");
                    $model_id  = "OS_BARCODE";
                    ;

                    break;
        case "1101":  include 'db_connection_graph.php';

                    $databaseInUse = "[YF_Clone].[dbo].[hanon_shaftrotor_table]";
                    //$name = "YF Line A";
                    $ID = "MINDEX";
                    $dateTime = "TIME_STAMP";
                    $param = "REJ_STATION = 0";
                    $paramRej = "REJECTED = 1";
                    $paramRejSta = "REJECTED = 1 AND REJ_STATION like '".$station."'";
                    $paramModel = ($model != "" ? $paramModel = "AND PART_NUMBER like '%$modelFAssy%' AND ": $paramModel = "AND");
                    $model_id  = "ROTOR_BARCODE";
                    ;

                    break;
        case "1102":  include 'db_connection_graph.php';

                    $databaseInUse = "[YF_Clone].[dbo].[hanon_screendiffuser_table]";
                    //$name = "YF Line A";
                    $ID = "MINDEX";
                    $dateTime = "TIME_STAMP";
                    $param = "REJ_STATION = 0";
                    $paramRej = "REJECTED = 1";
                    $paramRejSta = "REJECTED = 1 AND REJ_STATION like '".$station."'";
                    $paramModel = ($model != "" ? $paramModel = "AND PART_NUMBER like '%$modelFAssy%' AND ": $paramModel = "AND");
                    $model_id  = "RH_BARCODE";
                    ;

                    break;
        case "1103":  include 'db_connection_graph.php';

                $databaseInUse = "[YF_Clone].[dbo].[hanon_statorhm_table]";
                //$name = "YF Line A";
                $ID = "MINDEX";
                $dateTime = "TIME_STAMP";
                $param = "REJ_STATION = 0";
                $paramRej = "REJECTED = 1";
                $paramRejSta = "REJECTED = 1 AND REJ_STATION like '".$station."'";
                $paramModel = ($model != "" ? $paramModel = "AND PART_NUMBER like '%$modelFAssy%' AND ": $paramModel = "AND");
                $model_id  = "MH_BARCODE";
                ;

                break;
				

		default: 	$error = "<script>alert('No records found!');</script>";
					exit($error);

        }
?>