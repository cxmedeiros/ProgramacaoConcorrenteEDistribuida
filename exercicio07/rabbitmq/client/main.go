package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ArrayInt [1000]int

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func measureTime(ch *amqp.Channel, q amqp.Queue, replyQueueName string, body []byte) {
	correlationId := uuid.New().String()

	// startTime := time.Now()
	err := ch.PublishWithContext(
		context.Background(),
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			ReplyTo:       replyQueueName,
			CorrelationId: correlationId,
		})
	failOnError(err, "Failed to publish a message")

	// return correlationId
	// startTime := time.Now()
	// time.Since(startTime)
}

func main() {
	ArrayInt = [1000]int{451, 2961, 239, 1527, 3280, 4728, 4531, 3146, 2715, 3978, 98, 506, 1513, 2815, 4344, 4687, 1595, 4933, 3394, 2485, 536, 185, 279, 4588, 1653, 1792, 1269, 4168, 4505, 3235, 2690, 3815, 4644, 2840, 2712, 2163, 1850, 3353, 3466, 4868, 4633, 649, 1503, 1425, 2602, 1996, 3293, 4672, 4990, 4986, 1480, 1697, 750, 3937, 4786, 1917, 1777, 1733, 3693, 3125, 2299, 1688, 161, 2751, 4576, 331, 2846, 3604, 1162, 2381, 491, 1569, 2291, 2958, 2605, 2899, 4567, 1785, 1253, 1135, 1385, 228, 369, 1475, 4176, 1492, 436, 2348, 2382, 4055, 4651, 4916, 140, 3540, 3041, 1268, 1311, 4926, 1542, 880, 695, 791, 2696, 4383, 881, 3276, 2332, 2369, 2032, 401, 2080, 4126, 3584, 2051, 3355, 1586, 3631, 2211, 2656, 1536, 2003, 1097, 3629, 4655, 3643, 1220, 321, 1415, 484, 3354, 4362, 2235, 1517, 2881, 3823, 3381, 2252, 138, 4918, 63, 590, 656, 3044, 512, 4614, 1013, 3832, 193, 2661, 666, 776, 2398, 1325, 1823, 3444, 3339, 148, 738, 3175, 978, 3285, 3113, 989, 1280, 2204, 4250, 450, 3137, 3068, 4622, 2968, 2114, 1733, 4304, 1196, 973, 3680, 1607, 4577, 192, 4761, 1278, 3894, 4543, 1468, 2057, 2653, 1099, 4803, 1623, 1030, 1155, 160, 2399, 4132, 1433, 1023, 4645, 3002, 2803, 3072, 3611, 2662, 4364, 3646, 37, 2432, 2498, 3535, 3785, 3623, 29, 1268, 2644, 3945, 4556, 3067, 4133, 3245, 2088, 4402, 671, 225, 739, 4781, 449, 448, 365, 4143, 1410, 4734, 4190, 3353, 3416, 3792, 2967, 3614, 297, 1290, 2003, 2585, 4783, 4862, 702, 3622, 2592, 1482, 1300, 222, 3714, 4548, 3470, 2883, 2016, 4221, 4874, 2984, 818, 2568, 4207, 882, 443, 1256, 3640, 2948, 756, 1603, 1897, 2186, 2175, 179, 1272, 3269, 3565, 3829, 2775, 2919, 4795, 4813, 2248, 4915, 4107, 1982, 2487, 1768, 1238, 1535, 1429, 2747, 3661, 2197, 783, 2073, 4655, 4728, 3723, 1266, 1661, 3395, 775, 1297, 544, 3959, 3547, 4107, 3276, 2106, 4509, 2084, 404, 2405, 3683, 1940, 2527, 498, 1327, 2016, 3598, 1275, 1412, 4757, 461, 3160, 4415, 4622, 4530, 3792, 3482, 4122, 4139, 340, 877, 573, 2238, 4447, 2498, 137, 3460, 1139, 465, 3282, 2707, 1146, 405, 4848, 4367, 1748, 4170, 4095, 3655, 1366, 189, 3491, 692, 263, 2652, 2052, 1451, 3073, 4166, 240, 1908, 3115, 1395, 191, 3134, 652, 328, 4610, 335, 2054, 3538, 1762, 2482, 747, 4838, 3746, 228, 2778, 3429, 2636, 4621, 3023, 635, 3416, 3519, 1272, 2037, 2743, 3895, 740, 4299, 455, 258, 2302, 506, 3968, 1588, 156, 3612, 3767, 1193, 4747, 3051, 3040, 2112, 2730, 401, 2417, 1456, 1271, 987, 2437, 3970, 507, 931, 1274, 581, 720, 2196, 2211, 4247, 4624, 3665, 341, 2512, 4317, 4563, 125, 4090, 66, 2490, 3320, 515, 4038, 4848, 4130, 4065, 4967, 2447, 4866, 3679, 4275, 4795, 1255, 4450, 2245, 3003, 1234, 1799, 3263, 4793, 2215, 3282, 4415, 4653, 4642, 2767, 7, 2964, 1453, 1522, 1764, 4552, 3291, 1219, 722, 2823, 268, 4426, 2654, 3218, 4041, 1755, 2054, 2307, 117, 4679, 1345, 4588, 190, 3821, 4180, 4857, 3206, 3263, 1204, 2348, 1253, 4150, 790, 1359, 3531, 3728, 2828, 408, 132, 2632, 3710, 429, 4874, 2217, 1916, 433, 2967, 4836, 634, 2036, 3746, 439, 1212, 1842, 375, 4976, 1539, 3944, 3151, 2075, 760, 565, 976, 2564, 1464, 3315, 3093, 3646, 146, 3851, 834, 2652, 3757, 3861, 3878, 1366, 49, 715, 4099, 3976, 4095, 4462, 4834, 327, 2717, 898, 1550, 1929, 1780, 2646, 1997, 4470, 4695, 893, 4659, 3135, 1135, 4309, 2936, 430, 4564, 625, 4158, 2379, 4858, 663, 309, 476, 336, 1761, 1010, 2701, 4265, 3858, 3108, 4581, 4564, 4138, 3665, 4317, 3088, 280, 2410, 2330, 1820, 3006, 578, 1770, 1604, 2608, 1305, 1215, 2649, 662, 200, 4912, 3229, 2669, 2155, 292, 1794, 503, 2996, 2043, 2095, 4620, 100, 3874, 4799, 2180, 3689, 2760, 1441, 4419, 4530, 343, 900, 2738, 4102, 3349, 4354, 4054, 3696, 3180, 1075, 750, 727, 3126, 2734, 1311, 3154, 2453, 3767, 4393, 813, 2994, 720, 2056, 156, 411, 2054, 1168, 2798, 3518, 2150, 3746, 543, 2321, 1274, 3642, 1300, 1446, 1344, 1840, 4803, 570, 903, 3706, 1001, 484, 3003, 150, 3934, 4034, 1625, 3090, 1470, 4520, 2710, 3197, 614, 4326, 4863, 3185, 201, 2388, 3439, 4747, 931, 1436, 151, 696, 2460, 4241, 1869, 2130, 4474, 1202, 4209, 3158, 4848, 3738, 3583, 3765, 980, 374, 4234, 1488, 4783, 1559, 3633, 4299, 1177, 37, 4083, 4407, 1058, 4078, 2402, 3974, 1337, 1222, 572, 4476, 4176, 3925, 3852, 1138, 2686, 3089, 4869, 886, 3899, 4820, 4008, 1286, 4322, 3768, 2162, 3383, 2186, 2774, 3355, 2349, 2051, 2395, 4116, 304, 153, 2627, 4750, 4362, 2765, 4930, 1381, 4713, 4011, 224, 2834, 2971, 1372, 3154, 4128, 4603, 957, 4376, 2677, 3095, 2142, 1425, 4254, 4972, 3129, 1232, 4068, 1939, 446, 621, 4225, 2631, 4046, 4776, 1164, 2998, 2950, 759, 1065, 910, 4387, 4348, 4294, 3688, 27, 1051, 2145, 3495, 2763, 4523, 1408, 1785, 893, 2831, 1218, 1336, 2524, 3920, 622, 898, 1381, 187, 1473, 714, 3923, 2761, 1741, 1553, 4229, 344, 3853, 2739, 3850, 2943, 865, 3499, 2936, 32, 4830, 3446, 2716, 3687, 1226, 3617, 2273, 4968, 1331, 3754, 2339, 4809, 3226, 4480, 4524, 232, 3408, 2458, 4688, 319, 2619, 4671, 4896, 1004, 2515, 2207, 2733, 2915, 431, 4068, 3539, 3160, 2161, 4156, 162, 3178, 222, 3881, 75, 350, 3655, 1324, 4976, 183, 4208, 21, 151, 1270, 2851, 3876, 4481, 395, 2956, 3591, 289, 3472, 154, 3386, 3891, 3219, 3257, 1967, 3264, 3936, 3938, 4156, 623, 4420, 4638, 769, 717, 4388, 2858, 2650, 1949, 4754, 746, 3111, 969, 2336, 3165, 3764, 4692, 410, 2536, 2893, 2262, 979, 4671, 232, 4956, 2043, 1464, 3058, 642, 3493, 4420, 399, 1264, 3193, 1494, 822, 4575, 685, 773, 2871, 906, 1629, 1633, 3569, 4732, 1547, 4798, 4854, 681, 4828, 1659, 3320, 1461, 3632, 2986, 4055, 822, 1338, 2899, 437, 1378, 1575, 2987, 4167, 4312, 3973, 535, 3900, 1131, 1275, 3063, 4125, 3504, 1591, 847, 3260, 2766, 1900, 2760, 213, 3934, 3445, 63, 2566, 4959, 3177, 134, 217, 1059, 437, 1606, 1429, 15, 1650, 512, 4582, 4972, 3000, 621, 828, 3115, 539, 4067, 1414, 1861, 240, 2268, 3599, 3197, 4773, 4193, 3400, 1657, 3690, 2604, 1062, 2084, 638, 1367, 4003}

	var Array [1000]int64
	for i := 0; i < len(ArrayInt); i++ {
		Array[i] = int64(i)
	}

	countStr := os.Getenv("COUNT")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Fatalf("Invalid COUNT value: %v", err)
	}

	timestamp := time.Now().Format("2006-01-02_15h-04m-05s")
	filePath := "temp/message_times_" + timestamp + ".csv"
	csvFile, err := os.Create(filePath)
	failOnError(err, "Failed to create CSV file")
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	err = writer.Write([]string{"Num Messages", "Average Time (s)"})
	failOnError(err, "Failed to write CSV header")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	replyQueue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a reply queue")

	msgs, err := ch.Consume(
		replyQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	var wg sync.WaitGroup
	wg.Add(count)

	responses := make(map[string]string)

	go func() {
		for d := range msgs {
			responses[d.CorrelationId] = string(d.Body)
			wg.Done()
		}
	}()

	body, err := json.Marshal(Array)
	failOnError(err, "Failed to encode array to JSON")

	startTime := time.Now()
	for i := 0; i < count; i++ {
		measureTime(ch, q, replyQueue.Name, body)
	}
	wg.Wait()
	totalTime := time.Since(startTime)
	fmt.Printf("%v", totalTime)

	averageTime := totalTime.Seconds() / float64(count)
	fmt.Printf("%v", averageTime)
	fmt.Printf("Average time for %d messages (after trimming): %.6f seconds\n", count, averageTime)

	writer.Write([]string{fmt.Sprintf("%d", count), fmt.Sprintf("%.6f", averageTime)})
	writer.Flush()
}