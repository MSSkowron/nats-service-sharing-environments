import asyncio
import nats
import struct
import random

devices = {
}


async def publish_message():
    nc = await nats.connect("nats://admin:admin@nats2:4222")
    js = nc.jetstream()
    await nc.publish("publishers.publisher", struct.pack('?', True))

    async def device_handler(msg):
        device = msg.data.decode()

        device_type, device_room, device_name, device_data = device.split(".")

        if device_type not in devices:
            devices[device_type] = {}

        if device_room not in devices[device_type]:
            devices[device_type][device_room] = {}

        if device_name not in devices[device_type][device_room]:
            devices[device_type][device_room][device_name] = {}

        if devices[device_type][device_room][device_name] == {}:
            devices[device_type][device_room][device_name] = {device_data}

        print(f"Added new device: {device}")
    
    await js.subscribe("publishers.device", cb=device_handler)
    try:
        while True:
            await asyncio.sleep(1)
            rand = random.random()
            if rand < 0.5:
                await js.publish("lights.change", struct.pack('?', True))
                print("Published message: 'on'")
            else:
                await js.publish("lights.change", struct.pack('?', False))
                print("Published message: 'off'")
    except KeyboardInterrupt:
        print("Publisher stopped.")

    await nc.close()


if __name__ == '__main__':
    print("START")
    asyncio.run(publish_message())
