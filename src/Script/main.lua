Main = {}

local ServerIdList = {}

function Main.OnMessage(ser_id, table, func, param)
    print("srcName = "..table.."."..func.."("..param..")")
end

function Main.Send(serverId, table, func, param)
    if ServerIdList[serverId] then
        Main.SendMsg(serverId, table, func, param)
    else
        Main.Panic("不存在serverId:"..serverId)
    end
end

function Main.AddServer(ser_id, name, desc)
    local tmp = {}
    tmp.name = name
    tmp.id = ser_id
    tmp.desc = desc
    ServerIdList[ser_id] = tmp
end