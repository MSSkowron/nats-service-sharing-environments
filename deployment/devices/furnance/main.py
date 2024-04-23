import sys
import asyncio
import nats
import struct
import os
is_off = True
fuel = 50


async def main(name, room):
    global is_off, fuel

    nc = await nats.connect("nats://admin:admin@nats3:4222")
    js = nc.jetstream()

    print("Furnance connected")

    async def furnance_handler(msg):
          global is_off, fuel
          turnOn = struct.unpack('?', msg.data)[0]
          if turnOn:
            if fuel > 10:
              fuel-=10
              is_off = False
              reply = f"Success, furnance {name} turned on. Current fuel level is [{fuel}]."
              await msg.respond(reply.encode("utf8"))
              # await js.publish(msg.reply, reply)
            else:
              reply = f"Current fuel level [{fuel}] in {name} is too low"
              await msg.respond(reply.encode("utf8"))
              # await js.publish(msg.reply, reply)
          else:
            is_off = True
            reply = f"Success, furnance {name} turned off."
            await msg.respond(reply.encode("utf8"))
            # await js.publish(msg.reply, reply)
          reply = f"Default"
          await msg.respond(reply.encode("utf8"))

    async def status_handler(msg):
        global is_off
        await js.publish("publishers.device", f"furnances.{name}.{room}.{not is_off}".encode())

    

    await nc.subscribe(f"furnances.{room}.{name}.change", cb=furnance_handler)
    await nc.subscribe(f"furnances.{room}.change", cb=furnance_handler)
    await nc.subscribe("furnances.change", cb=furnance_handler)

    
    await js.subscribe("publishers.publisher", cb=status_handler)   
    await js.publish("publishers.device", f"furnances.{room}.{name}.{is_off}".encode())

    try:
        while True:
            await asyncio.sleep(1)
            fuel+=5
    except KeyboardInterrupt:
        print("Subscriber stopped.")

    await nc.close()


if __name__ == '__main__':
    name = os.environ.get('FURNANCE_NAME')
    room = os.environ.get('FURNANCE_ROOM')
    asyncio.run(main(name,room))
