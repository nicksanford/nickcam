# Nick Cam

### Build:
```
make
```

### Example Config:
```json
{
  "components": [
    {
      "name": "nickcam",
      "namespace": "rdk",
      "type": "camera",
      "model": "ncs:camera:nickcam",
      "attributes": {
        "big": true,
        "color": "green",
        "image_type": "jpeg"
      },
    }
  ],
  "modules": [
    {
      "type": "local",
      "name": "ncs-nickcam",
      "executable_path": "/home/user/nickcam"
    }
  ]
}
```
