package message
//  主要负责消息与二进制之间的转换
import (
	"encoding/binary"
	"math"
)

const(
	MESSAGEBASE       byte = 0
	BYTE              byte = 1
	INT16             byte = 2
	INT32             byte = 3
	INT64             byte = 4
	BOOL              byte = 5
	FLOAT32           byte = 6
	FLOAT64           byte = 7
	STRING            byte = 8
	BINARY            byte = 9
	USE_LITTLE_ENDIAN      = true
)

func GetMsgId(propertyList []interface{}) uint32{
	return uint32(propertyList[1].(int32))
}

func GetServerId(propertyList []interface{}) byte{
	return byte(propertyList[2].(byte))
}

func WriteByte(PropertyList []interface{}) ([]byte, error){
	if USE_LITTLE_ENDIAN {
		var littleEndian = binary.LittleEndian
		var btsLength = make([]byte, 2)
		var data = make([]byte, 0)
		var length = uint16(len(PropertyList))
		//fmt.Println("length =", length)
		littleEndian.PutUint16(btsLength, length)
		data = append(data, btsLength...)
		//fmt.Println(PropertyList)
		for index := 0; index < len(PropertyList); index ++{

			switch v := PropertyList[index].(type){
			case byte:
				data = append(data, BYTE)
				data = append(data, v)
			case int16:
				data = append(data, INT16)
				var bytes = make([]byte, 2)
				littleEndian.PutUint16(bytes, uint16(v))
				data = append(data, bytes...)
			case int32:
				data = append(data, INT32)
				var bytes = make([]byte, 4)
				littleEndian.PutUint32(bytes, uint32(v))
				data = append(data, bytes...)
			case uint32:
				data = append(data, INT32)
				var bytes = make([]byte, 4)
				littleEndian.PutUint32(bytes, uint32(v))
				data = append(data, bytes...)
			case int64:
				data = append(data, INT64)
				var bytes = make([]byte, 8)
				littleEndian.PutUint64(bytes, uint64(v))
				data = append(data, bytes...)
			case bool:
				data = append(data, BOOL)
				if v == true {
					var bt byte = 1
					data = append(data, bt)
				} else {
					var bt byte = 0
					data = append(data, bt)
				}
			case float32:
				data = append(data, FLOAT32)
				bits := math.Float32bits(v)
				bytes := make([]byte, 4)
				littleEndian.PutUint32(bytes, bits)
				data = append(data, bytes...)
			case float64:
				data = append(data, FLOAT64)
				bits := math.Float64bits(v)
				bytes := make([]byte, 8)
				littleEndian.PutUint64(bytes, bits)
				data = append(data, bytes...)
			case string:
				data = append(data, STRING)
				var bytes = []byte(v)
				var length = int16(len(bytes))
				var bytesLength = make([]byte, 2)
				littleEndian.PutUint16(bytesLength, uint16(length))
				data = append(data, bytesLength...)
				data = append(data, bytes...)
			case []byte:
				data = append(data, BINARY)
				var length = len(v)
				var bytes = make([]byte, 4)
				littleEndian.PutUint32(bytes, uint32(length))
				data = append(data, bytes...)
				data = append(data, v...)
			case IMessage:
				data = append(data, MESSAGEBASE)
				var proertyList = v.GetPropertyList()
				bytes,err := WriteByte(proertyList)
				if err != nil {
					return nil, err
				}
				var length = int32(len(bytes)) // 写入长度
				var bytes1 = make([]byte, 4)
				littleEndian.PutUint32(bytes1, uint32(length))
				data = append(data, bytes1...)
				data = append(data, bytes...)
			}
		}

		return data, nil
	} else {
		var bigEndian = binary.BigEndian
		var btsLength = make([]byte, 2)
		var data = make([]byte, 0)
		var length = uint16(len(PropertyList))
		//fmt.Println("length =", length)
		bigEndian.PutUint16(btsLength, length)
		data = append(data, btsLength...)

		for index := 0; index < len(PropertyList); index ++{
			switch v := PropertyList[index].(type){
			case byte:
				data = append(data, BYTE)
				data = append(data, v)
			case int16:
				data = append(data, INT16)
				var bytes = make([]byte, 2)
				bigEndian.PutUint16(bytes, uint16(v))
				data = append(data, bytes...)
			case int32:
				data = append(data, INT32)
				var bytes = make([]byte, 4)
				bigEndian.PutUint32(bytes, uint32(v))
				data = append(data, bytes...)
			case uint32:
				data = append(data, INT32)
				var bytes = make([]byte, 4)
				bigEndian.PutUint32(bytes, uint32(v))
				data = append(data, bytes...)
			case int64:
				data = append(data, INT64)
				var bytes = make([]byte, 8)
				bigEndian.PutUint64(bytes, uint64(v))
				data = append(data, bytes...)
			case bool:
				data = append(data, BOOL)
				if v == true {
					var bt byte = 1
					data = append(data, bt)
				} else {
					var bt byte = 0
					data = append(data, bt)
				}
			case float32:
				data = append(data, FLOAT32)
				bits := math.Float32bits(v)
				bytes := make([]byte, 4)
				bigEndian.PutUint32(bytes, bits)
				data = append(data, bytes...)
			case float64:
				data = append(data, FLOAT64)
				bits := math.Float64bits(v)
				bytes := make([]byte, 8)
				bigEndian.PutUint64(bytes, bits)
				data = append(data, bytes...)
			case string:
				data = append(data, STRING)
				var bytes = []byte(v)
				var length = int16(len(bytes))
				var bytesLength = make([]byte, 2)
				bigEndian.PutUint16(bytesLength, uint16(length))
				data = append(data, bytesLength...)
				data = append(data, bytes...)
			case []byte:
				data = append(data, BINARY)
				var length = len(v)
				var bytes = make([]byte, 4)
				bigEndian.PutUint32(bytes, uint32(length))
				data = append(data, bytes...)
				data = append(data, v...)
			case IMessage:
				data = append(data, MESSAGEBASE)
				var proertyList = v.GetPropertyList()
				bytes,err := WriteByte(proertyList)
				if err != nil {
					return nil, err
				}
				var length = int32(len(bytes)) // 写入长度
				var bytes1 = make([]byte, 4)
				bigEndian.PutUint32(bytes1, uint32(length))
				data = append(data, bytes1...)
				data = append(data, bytes...)
			}
		}
		return data, nil
	}
}

func GetPropertyList(data []byte) []interface{}{
	var PropertyList = make([]interface{}, 0)
	if(USE_LITTLE_ENDIAN){
		var littleEndian = binary.LittleEndian
		var length = len(data)
		//fmt.Println("length(PropertyList) =", length)
		var bytes = data[:2]
		PropertyList = append(PropertyList, int16(littleEndian.Uint16(bytes)))
		//fmt.Print("len:", int16(littleEndian.Uint16(bytes)), "\n")
		//debug.PrintStack()
		for i := 2; i < length;{
			var btType = data[i]
			i++
			//fmt.Println("btType =", btType)
			switch btType {
			case BYTE:
				PropertyList = append(PropertyList, data[i])
				//fmt.Print("BYTE:", data[i])
				i++
			case INT16:
				var bytes = data[i : i + 2]
				PropertyList = append(PropertyList, int16(littleEndian.Uint16(bytes)))
				i += 2
			case INT32:
				var bytes = data[i : i + 4]
				PropertyList = append(PropertyList, int32(littleEndian.Uint32(bytes)))
				//fmt.Print("INT32:", int32(littleEndian.Uint32(bytes)))
				i += 4
			case INT64:
				var bytes = data[i : i + 8]
				PropertyList = append(PropertyList, int64(littleEndian.Uint64(bytes)))
				i += 8
			case BOOL:
				var bt = data[i]
				i++
				PropertyList = append(PropertyList, bt == 1)
				//fmt.Print("INT32:", bt == 1, ", i =", i, ", data =", data[i:])
			case FLOAT32:
				var bytes = data[i : i + 4]
				var vUint32 = littleEndian.Uint32(bytes)
				PropertyList = append(PropertyList, math.Float32frombits(vUint32))
				i += 4
			case FLOAT64:
				var bytes = data[i : i + 8]
				var vUint64 = littleEndian.Uint64(bytes)
				PropertyList = append(PropertyList, math.Float64frombits(vUint64))
				i += 8
			case STRING:
				var bytes = data[i : i + 2]
				i += 2
				var length = int(littleEndian.Uint16(bytes))
				bytes = data[i : i + length]
				i += length
				PropertyList = append(PropertyList, string(bytes))
			case BINARY:
				var bytes = data[i : i + 4]
				i += 4
				var length = int(littleEndian.Uint32(bytes))
				bytes = data[i : i + length]
				i += length
				PropertyList = append(PropertyList, bytes)
			case MESSAGEBASE:
				var bytes = data[i : i + 4]
				i += 4
				var length = int(int32(littleEndian.Uint32(bytes)))
				bytes = data[i : i + length]
				i += length
				var propertyList = GetPropertyList(bytes)
				var message = GetMessage(propertyList)
				PropertyList = append(PropertyList, message)
			}
		}
	} else {
		var bigEndian = binary.BigEndian
		var length = len(data)
		//fmt.Println("length(PropertyList) =", length)
		var bytes = data[:2]
		PropertyList = append(PropertyList, int16(bigEndian.Uint16(bytes)))
		for i := 2; i < length;{
			var btType = data[i]
			i++
			//fmt.Println("btType =", btType)
			switch btType {
			case BYTE:
				PropertyList = append(PropertyList, data[i])
				i++
			case INT16:
				var bytes = data[i : i + 2]
				PropertyList = append(PropertyList, int16(bigEndian.Uint16(bytes)))
				i += 2
			case INT32:
				var bytes = data[i : i + 4]
				PropertyList = append(PropertyList, int32(bigEndian.Uint32(bytes)))
				i += 4
			case INT64:
				var bytes = data[i : i + 8]
				PropertyList = append(PropertyList, int64(bigEndian.Uint64(bytes)))
				i += 8
			case BOOL:
				var bt = data[i]
				i++
				PropertyList = append(PropertyList, bt == 1)
			case FLOAT32:
				var bytes = data[i : i + 4]
				var vUint32 = bigEndian.Uint32(bytes)
				PropertyList = append(PropertyList, math.Float32frombits(vUint32))
				i += 4
			case FLOAT64:
				var bytes = data[i : i + 8]
				var vUint64 = bigEndian.Uint64(bytes)
				PropertyList = append(PropertyList, math.Float64frombits(vUint64))
				i += 8
			case STRING:
				var bytes = data[i : i + 2]
				i += 2
				var length = int(bigEndian.Uint16(bytes))
				bytes = data[i : i + length]
				i += length
				PropertyList = append(PropertyList, string(bytes))
			case BINARY:
				var bytes = data[i : i + 4]
				i += 4
				var length = int(bigEndian.Uint32(bytes))
				bytes = data[i : i + length]
				i += length
				PropertyList = append(PropertyList, bytes)
			case MESSAGEBASE:
				var bytes = data[i : i + 4]
				i += 4
				var length = int(int32(bigEndian.Uint32(bytes)))
				bytes = data[i : i + length]
				i += length
				var propertyList = GetPropertyList(bytes)
				var message = GetMessage(propertyList)
				PropertyList = append(PropertyList, message)
			}
		}
	}
	return PropertyList
}