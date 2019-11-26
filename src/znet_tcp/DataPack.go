package znet_tcp

import (
	"bytes"
	"message"
	"fmt"
	"encoding/binary"
)

//封包拆包类实例，暂时不需要成员
type DataPack struct {
	Seq uint16	// 消息序号
	Ack uint16	// 消息确认号
	Body []byte	// 消息的内容
	Data []byte	// 发送的内容
	IsBuff bool	// 是否使用的带缓冲的chan
}

//拆包
func NewUnDataPack(data []byte) *DataPack {
	if len(data) >= 4 {
		if message.USE_LITTLE_ENDIAN {
			var littleEndian = binary.LittleEndian
			var bts_seq = data[0 : 2]
			var bts_ack = data[2 : 4]
			var bts_data = data[4 :]
			return &DataPack{
				Seq:littleEndian.Uint16(bts_seq),
				Ack:littleEndian.Uint16(bts_ack),
				Body:bts_data,
				Data:data,
			}
		} else {
			var bigEndian = binary.LittleEndian
			var bts_seq = data[0 : 2]
			var bts_ack = data[2 : 4]
			var bts_data = data[4 :]
			return &DataPack{
				Seq:bigEndian.Uint16(bts_seq),
				Ack:bigEndian.Uint16(bts_ack),
				Body:bts_data,
				Data:data,
			}
		}
	} else {
		fmt.Errorf("消息格式错误！")
		return nil
	}
}

//封包
func NewDataPack(data []byte, seq uint16, ack uint16, isBuff bool) *DataPack {
	if message.USE_LITTLE_ENDIAN {	// 使用小端

		var littleEndian = binary.LittleEndian

		var data_pack = make([]byte, 0)
		PH_LEN := int32(len(data)) + 4
		var bts_len = make([]byte, 4)
		littleEndian.PutUint32(bts_len, uint32(PH_LEN))
		data_pack = append(data_pack, bts_len...)

		// 写入序号
		var bts_seq = make([]byte, 2)
		littleEndian.PutUint16(bts_seq, seq)
		data_pack = append(data_pack, bts_seq...)
		// 写入消息确认号
		var bts_ack = make([]byte, 2)
		littleEndian.PutUint16(bts_ack, ack)
		data_pack = append(data_pack, bts_ack...)
		// 写入消息体
		data_pack = append(data_pack, data...)
		return &DataPack{
			Seq:seq,
			Ack:ack,
			Body:data,
			Data:data_pack,
			IsBuff:isBuff,
		}
	} else {
		var bigEndian = binary.LittleEndian

		var data_pack = make([]byte, 0)
		PH_LEN := int32(len(data)) + 4
		var bts_len = make([]byte, 4)
		bigEndian.PutUint32(bts_len, uint32(PH_LEN))
		data_pack = append(data_pack, bts_len...)

		// 写入序号
		var bts_seq = make([]byte, 2)
		bigEndian.PutUint16(bts_seq, seq)
		data_pack = append(data_pack, bts_seq...)
		// 写入消息确认号
		var bts_ack = make([]byte, 2)
		bigEndian.PutUint16(bts_ack, ack)
		data_pack = append(data_pack, bts_ack...)
		// 写入消息体
		data_pack = append(data_pack, data...)
		return &DataPack{
			Seq:seq,
			Ack:ack,
			Body:data,
			Data:data_pack,
			IsBuff:isBuff,
		}
	}
}

//获取包头长度方法
func(dp *DataPack) GetHeadLen() int {
	return 4
}

//封包方法(压缩数据)
func(dp *DataPack) Pack(msg message.IMessage)([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	/*
	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil ,err
	}
	//*/
	return dataBuff.Bytes(), nil
}
//拆包方法(解压数据)
func(dp *DataPack) Unpack(binaryData []byte)(message.IMessage, error) {
	/*
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}


	//判断dataLen的长度是否超出我们允许的最大包长度
	if (utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize) {
		return nil, errors.New("Too large msg data recieved")
	}
	//*/
	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return nil, nil
}
