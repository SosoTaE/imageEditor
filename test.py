import requests
import os
import base64
import time

body = {"objects":[
    {
        "images": [],
        "functions": [{
            "name": "crop",
            "parameters": {
                "coordinate":[0,0,1920,1080]
            }
        },``
        {
            "name": "resize",
            "parameters": {
                "size":[3840,2160]
            }
        }
        ]
    }
]}


folder = "D:\\pythonProject2\\test_folder"

urls = [os.path.join(folder,url) for url in os.listdir(folder)]

for each in urls:
    with open(each, "rb") as file:
        data = file.read()

        base64_string = base64.b64encode(data).decode("utf-8")

        body["objects"][0]["images"].append(base64_string)

start = time.time()

# print(body)

response = requests.post("http://localhost:9000/api/v1/editor", json=body)

if not response.ok:
   exit()

data = response.json()

counter = 0
for index, each in enumerate(data.get("array")):
   for img in each:
        byte = base64.b64decode(img)
        with open(f"{counter}.jpg", "wb") as file:
            file.write(byte)
            counter +=1
        
end = time.time()


print(f"execution time:{end-start} image length:{counter}")