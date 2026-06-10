package seed

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"qc-system/internal/model"
	"qc-system/internal/service/aql"
	"qc-system/internal/store"
)

var (
	materialTypes = []model.MaterialType{
		model.MaterialMetal,
		model.MaterialPlastic,
		model.MaterialElectronic,
		model.MaterialChemical,
		model.MaterialPaper,
	}

	materialPrefixes = map[model.MaterialType][]string{
		model.MaterialMetal:      {"AL", "ST", "CU", "ZN", "TI"},
		model.MaterialPlastic:    {"PP", "PE", "PVC", "ABS", "PC"},
		model.MaterialElectronic: {"R", "C", "D", "T", "IC"},
		model.MaterialChemical:   {"CH", "SOL", "POW", "LIQ", "GEL"},
		model.MaterialPaper:      {"PAP", "BOX", "BAG", "LAB", "TAG"},
	}

	materialNames = map[model.MaterialType][]string{
		model.MaterialMetal:      {"铝合金板材", "不锈钢管", "铜棒", "锌合金", "钛合金"},
		model.MaterialPlastic:    {"聚丙烯颗粒", "聚乙烯膜", "PVC管", "ABS外壳", "PC镜片"},
		model.MaterialElectronic: {"电阻", "电容", "二极管", "三极管", "集成电路"},
		model.MaterialChemical:   {"化学试剂", "溶剂", "粉末", "液体", "凝胶"},
		model.MaterialPaper:      {"包装纸", "纸箱", "纸袋", "标签", "吊牌"},
	}

	inspectionTypes = []model.InspectionType{
		model.TypeIQC,
		model.TypeIPQC,
		model.TypeOQC,
		model.TypeOBC,
	}

	judgments = []model.JudgmentResult{
		model.JudgmentPass,
		model.JudgmentPass,
		model.JudgmentPass,
		model.JudgmentPass,
		model.JudgmentFail,
		model.JudgmentReject,
	}

	dispositions = []model.DispositionType{
		model.DispositionConcession,
		model.DispositionRework,
		model.DispositionScrap,
		model.DispositionReturn,
	}

	shifts = []string{"早班", "中班", "晚班"}
	productionLines = []string{"A线", "B线", "C线", "D线"}
	equipments = []string{"设备01", "设备02", "设备03", "设备04", "设备05"}
	workers = []string{"张三", "李四", "王五", "赵六", "钱七", "孙八", "周九", "吴十"}

	inspectorNames = []string{"检验员A", "检验员B", "检验员C", "检验员D", "检验员E", "检验员F", "检验员G", "检验员H"}
)

func Seed() {
	rand.Seed(time.Now().UnixNano())

	suppliers := seedSuppliers(50)
	for _, s := range suppliers {
		store.GlobalStore.SaveSupplier(s)
	}

	materials := seedMaterials(200, suppliers)
	for _, m := range materials {
		store.GlobalStore.SaveMaterial(m)
	}

	inspectors := seedInspectors(8)
	for _, i := range inspectors {
		store.GlobalStore.SaveInspector(i)
	}

	seedStandards(materials)

	lots := seedLots(500, materials, suppliers)
	for _, l := range lots {
		store.GlobalStore.SaveLot(l)
	}

	seedRecords(5000, lots, materials, suppliers, inspectors)
}

func seedSuppliers(count int) []*model.Supplier {
	suppliers := make([]*model.Supplier, count)
	supplierNames := []string{
		"华为科技", "中兴通讯", "比亚迪", "宁德时代", "小米科技",
		"富士康", "立讯精密", "歌尔股份", "蓝思科技", "立讯精密",
		"瑞声科技", "舜宇光学", "京东方", "TCL科技", "三安光电",
		"韦尔股份", "兆易创新", "紫光国微", "中芯国际", "长电科技",
		"华天科技", "通富微电", "北方华创", "中微公司", "北方华创",
		"汇顶科技", "卓胜微", "圣邦股份", "韦尔股份", "北京君正",
		"兆驰股份", "木林森", "国星光电", "鸿利智汇", "聚飞光电",
		"瑞丰光电", "利亚德", "洲明科技", "艾比森", "雷曼光电",
		"长信科技", "凯盛科技", "东旭光电", "彩虹股份", "维信诺",
		"深天马", "京东方A", "TCL科技", "华星光电", "天马微电子",
	}

	for i := 0; i < count; i++ {
		suppliers[i] = &model.Supplier{
			ID:        uuid.New().String(),
			Code:      fmt.Sprintf("SUP%04d", i+1),
			Name:      supplierNames[i%len(supplierNames)] + fmt.Sprintf("-%d", i/len(supplierNames)+1),
			Contact:   fmt.Sprintf("联系人%d", i+1),
			Phone:     fmt.Sprintf("138%08d", i+1),
			Address:   fmt.Sprintf("广东省深圳市南山区科技园%d栋", i+1),
			Rating:    60 + rand.Float64()*35,
			Status:    "normal",
			CreatedAt: time.Now().AddDate(0, -rand.Intn(12), -rand.Intn(30)),
			UpdatedAt: time.Now(),
		}
	}
	return suppliers
}

func seedMaterials(count int, suppliers []*model.Supplier) []*model.Material {
	materials := make([]*model.Material, count)

	for i := 0; i < count; i++ {
		matType := materialTypes[i%len(materialTypes)]
		prefixes := materialPrefixes[matType]
		names := materialNames[matType]

		materials[i] = &model.Material{
			ID:            uuid.New().String(),
			Code:          fmt.Sprintf("%s%04d", prefixes[i%len(prefixes)], i+1),
			Name:          names[i%len(names)],
			Type:          matType,
			Specification: fmt.Sprintf("规格%d", i+1),
			Unit:          "pcs",
			SupplierID:    suppliers[i%len(suppliers)].ID,
			CreatedAt:     time.Now().AddDate(0, -rand.Intn(12), -rand.Intn(30)),
			UpdatedAt:     time.Now(),
		}
	}
	return materials
}

func seedInspectors(count int) []*model.Inspector {
	inspectors := make([]*model.Inspector, count)
	roles := []string{"inspector", "inspector", "inspector", "qa_manager", "inspector", "inspector", "director", "inspector"}

	for i := 0; i < count; i++ {
		inspectors[i] = &model.Inspector{
			ID:         uuid.New().String(),
			Name:       inspectorNames[i],
			EmployeeNo: fmt.Sprintf("INS%04d", i+1),
			Role:       roles[i],
			Processes:  []string{fmt.Sprintf("工序%d", i+1), fmt.Sprintf("工序%d", i+2)},
			Status:     "active",
			Phone:      fmt.Sprintf("139%08d", i+100),
			Email:      fmt.Sprintf("inspector%d@company.com", i+1),
			CreatedAt:  time.Now().AddDate(0, -rand.Intn(12), -rand.Intn(30)),
			UpdatedAt:  time.Now(),
		}
	}
	return inspectors
}

func seedStandards(materials []*model.Material) {
	for _, mat := range materials {
		defaultItems := getDefaultItems(mat.Type)

		for _, t := range inspectionTypes {
			std := &model.Standard{
				ID:               uuid.New().String(),
				MaterialID:       mat.ID,
				Type:             t,
				InspectionLevels: defaultItems,
				AQL:              1.0,
				InspectionLevel:  "II",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}
			store.GlobalStore.SaveStandard(std)
		}
	}
}

func getDefaultItems(materialType model.MaterialType) []model.InspectionItem {
	items := []model.InspectionItem{
		{ID: uuid.New().String(), Name: "外观", Standard: "无明显划痕、变形", IsNumeric: false, SPCEnabled: false},
		{ID: uuid.New().String(), Name: "尺寸", Standard: "10.0±0.1mm", Unit: "mm", MinValue: 9.9, MaxValue: 10.1, IsNumeric: true, SPCEnabled: true},
		{ID: uuid.New().String(), Name: "性能", Standard: "符合规格要求", IsNumeric: false, SPCEnabled: false},
	}

	switch materialType {
	case model.MaterialMetal:
		items = append(items,
			model.InspectionItem{ID: uuid.New().String(), Name: "硬度", Standard: "HRC 20-30", Unit: "HRC", MinValue: 20, MaxValue: 30, IsNumeric: true, SPCEnabled: true},
			model.InspectionItem{ID: uuid.New().String(), Name: "化学成分", Standard: "符合标准", IsNumeric: false, SPCEnabled: false},
		)
	case model.MaterialPlastic:
		items = append(items,
			model.InspectionItem{ID: uuid.New().String(), Name: "颜色", Standard: "无色差", IsNumeric: false, SPCEnabled: false},
			model.InspectionItem{ID: uuid.New().String(), Name: "强度", Standard: "≥50MPa", Unit: "MPa", MinValue: 50, MaxValue: 100, IsNumeric: true, SPCEnabled: true},
		)
	case model.MaterialElectronic:
		items = append(items,
			model.InspectionItem{ID: uuid.New().String(), Name: "电阻", Standard: "100±5Ω", Unit: "Ω", MinValue: 95, MaxValue: 105, IsNumeric: true, SPCEnabled: true},
			model.InspectionItem{ID: uuid.New().String(), Name: "耐压", Standard: "≥500V", Unit: "V", MinValue: 500, MaxValue: 1000, IsNumeric: true, SPCEnabled: false},
		)
	case model.MaterialChemical:
		items = append(items,
			model.InspectionItem{ID: uuid.New().String(), Name: "纯度", Standard: "≥99.5%", Unit: "%", MinValue: 99.5, MaxValue: 100, IsNumeric: true, SPCEnabled: true},
			model.InspectionItem{ID: uuid.New().String(), Name: "水分", Standard: "≤0.1%", Unit: "%", MinValue: 0, MaxValue: 0.1, IsNumeric: true, SPCEnabled: false},
		)
	case model.MaterialPaper:
		items = append(items,
			model.InspectionItem{ID: uuid.New().String(), Name: "克重", Standard: "200±5g/m²", Unit: "g/m²", MinValue: 195, MaxValue: 205, IsNumeric: true, SPCEnabled: true},
			model.InspectionItem{ID: uuid.New().String(), Name: "厚度", Standard: "0.2±0.02mm", Unit: "mm", MinValue: 0.18, MaxValue: 0.22, IsNumeric: true, SPCEnabled: false},
		)
	}

	return items
}

func seedLots(count int, materials []*model.Material, suppliers []*model.Supplier) []*model.Lot {
	lots := make([]*model.Lot, count)

	lotCounters := make(map[string]int)

	for i := 0; i < count; i++ {
		mat := materials[i%len(materials)]
		sup := suppliers[i%len(suppliers)]

		yearMonth := time.Now().AddDate(0, -rand.Intn(6), -rand.Intn(30)).Format("200601")
		key := mat.Code + "-" + yearMonth
		lotCounters[key]++
		seq := lotCounters[key]

		lots[i] = &model.Lot{
			ID:             uuid.New().String(),
			LotNo:          fmt.Sprintf("%s-%s-%04d", mat.Code, yearMonth, seq),
			MaterialID:     mat.ID,
			Quantity:       100 + rand.Intn(5000),
			ProductionDate: time.Now().AddDate(0, -rand.Intn(6), -rand.Intn(30)),
			Shift:          shifts[rand.Intn(len(shifts))],
			Worker:         workers[rand.Intn(len(workers))],
			ProductionLine: productionLines[rand.Intn(len(productionLines))],
			Equipment:      equipments[rand.Intn(len(equipments))],
			RawMaterialLot: fmt.Sprintf("RAW%06d", rand.Intn(10000)),
			SupplierID:     sup.ID,
			Status:         "completed",
			CreatedAt:      time.Now().AddDate(0, -rand.Intn(6), -rand.Intn(30)),
		}
	}
	return lots
}

func seedRecords(count int, lots []*model.Lot, materials []*model.Material, suppliers []*model.Supplier, inspectors []*model.Inspector) {
	statuses := []model.RecordStatus{
		model.StatusClosed,
		model.StatusClosed,
		model.StatusClosed,
		model.StatusApproved,
		model.StatusSubmitted,
		model.StatusDraft,
	}

	for i := 0; i < count; i++ {
		lot := lots[i%len(lots)]
		mat, _ := store.GlobalStore.GetMaterial(lot.MaterialID)
		sup, _ := store.GlobalStore.GetSupplier(lot.SupplierID)
		inspector := inspectors[i%len(inspectors)]
		inspectionType := inspectionTypes[i%len(inspectionTypes)]

		status := statuses[rand.Intn(len(statuses))]
		judgment := judgments[rand.Intn(len(judgments))]

		std, ok := store.GlobalStore.GetStandardByMaterial(mat.ID, inspectionType)
		if !ok {
			continue
		}

		items := make([]model.InspectionItemResult, 0, len(std.InspectionLevels))
		defectCount := 0

		for _, item := range std.InspectionLevels {
			itemJudgment := model.JudgmentPass
			actualValue := "合格"
			numericValue := 0.0

			if item.IsNumeric {
				baseValue := (item.MinValue + item.MaxValue) / 2
				rangeVal := (item.MaxValue - item.MinValue) / 2
				numericValue = baseValue + (rand.Float64()-0.5)*rangeVal*1.5
				actualValue = fmt.Sprintf("%.3f", numericValue)

				if numericValue < item.MinValue || numericValue > item.MaxValue {
					itemJudgment = model.JudgmentFail
					defectCount++
				}
			} else if rand.Float64() < 0.1 {
				itemJudgment = model.JudgmentFail
				actualValue = "有瑕疵"
				defectCount++
			}

			items = append(items, model.InspectionItemResult{
				ItemID:       item.ID,
				ItemName:     item.Name,
				Standard:     item.Standard,
				ActualValue:  actualValue,
				NumericValue: numericValue,
				IsNumeric:    item.IsNumeric,
				Judgment:     itemJudgment,
				DefectDesc:   "",
			})
		}

		sampleResult, _ := aql.CalculateSample(lot.Quantity, std.AQL, std.InspectionLevel)
		sampleSize := lot.Quantity
		if sampleResult != nil {
			sampleSize = sampleResult.SampleSize
		}

		isFullInspection := inspectionType == model.TypeOQC && lot.Quantity <= 1000

		var disposition model.DispositionType
		if judgment != model.JudgmentPass {
			disposition = dispositions[rand.Intn(len(dispositions))]
		}

		createdAt := time.Now().AddDate(0, -rand.Intn(6), -rand.Intn(30))
		submittedAt := createdAt.Add(time.Hour * time.Duration(rand.Intn(24)))
		approvedAt := submittedAt.Add(time.Hour * time.Duration(rand.Intn(48)))
		closedAt := approvedAt.Add(time.Hour * time.Duration(rand.Intn(24)))

		record := &model.InspectionRecord{
			ID:               uuid.New().String(),
			LotID:            lot.ID,
			LotNo:            lot.LotNo,
			MaterialID:       mat.ID,
			SupplierID:       sup.ID,
			Type:             inspectionType,
			InspectorID:      inspector.ID,
			InspectorName:    inspector.Name,
			Status:           status,
			TotalSampleSize:  sampleSize,
			DefectCount:      defectCount,
			FinalJudgment:    judgment,
			Disposition:      disposition,
			Items:            items,
			Version:          1,
			Remarks:          []model.Remark{},
			IsFullInspection: isFullInspection,
			CreatedAt:        createdAt,
		}

		if status != model.StatusDraft {
			record.SubmittedAt = &submittedAt
		}
		if status == model.StatusApproved || status == model.StatusClosed {
			record.ApprovedAt = &approvedAt
		}
		if status == model.StatusClosed {
			record.ClosedAt = &closedAt
		}

		store.GlobalStore.SaveRecord(record)
	}
}
