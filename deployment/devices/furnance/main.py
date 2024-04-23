import sys
import asyncio
import nats
import struct

is_off = True
fuel = 50


async def main(name, room):
    global is_off

    nc = await nats.connect("nats://admin:admin@nats3:4222")
    js = nc.jetstream()

    print("Furnance connected")

    async def furnance_handler(msg):
          global is_off, fuel
          turnOn = struct.unpack('?', msg.data)[0]
          if turnOn and not is_off:
            if fuel > 10:
              fuel-=10
              is_off = False
              reply = f"Success, furnance {name} turned on. Current fuel level is [{fuel}]."
              await msg.respond(reply.encode("utf8"))
            else:
              reply = f"Current fuel level [{fuel}] in {name} is too low"
              await msg.respond(reply.encode("utf8"))
          elif not turnOn:
            is_off = True
            reply = f"Success, furnance {name} turned off."
            await msg.respond(reply.encode("utf8"))

          print(f"Received message: {turnOn}")

    async def status_handler(msg):
        global is_off
        await js.publish("publishers.device", f"furnances.{name}.{room}.{not is_off}".encode())

    

    await js.subscribe(f"furnances.{room}.{name}.change", cb=furnance_handler)
    await js.subscribe(f"furnances.{room}.change", cb=furnance_handler)
    await js.subscribe("furnances.change", cb=furnance_handler)

    
    await js.subscribe("publishers.publisher", cb=status_handler)   
    await js.publish("publishers.device", f"furnances.{room}.{name}.{is_off}".encode())

    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        print("Subscriber stopped.")

    await nc.close()


if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("Usage: python furnance.py <furnance_name> <furnance_room>")
        sys.exit(1)
    name = sys.argv[1]
    room = sys.argv[1]
    asyncio.run(main(name,room))
