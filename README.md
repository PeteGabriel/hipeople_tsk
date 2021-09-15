Web API that allows users to upload images and access them later. 



## API routes

Documentation on the entry points of this web API.

### Upload an image

Upload an image. An identifier is returned to be used later on when retrieving the image. 
Size of image is limited to 1MB.


* **URL**

  /api/image

* **Method:**

  `POST`

* **URL Params**

  None

* **Data Params**

  * **Multipart Form:** upload_file <br />
    **Content:** a file (up to 1MB)

* **Headers**

  * **Content-Type:** 'multipart/form-data;'

* **Success Response:**

    * **Code:** 201 <br />
      **Content-Type:** application/json <br />
      **Example:** `{
      "image_id": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
      }`

* **Error Response:**

    * **Code:** 404 NOT FOUND <br />
      **Content-Type:** application/problem+json <br />
      **Example:** 
      ``` 
      {
        "title": "file received is invalid",
        "status": 400,
        "detail": "File type not valid for upload. File received must be an image.",
        "instance": "/api/image"
      }
      ```

    * **Code:** 500 SERVER ERROR <br />
      **Content-Type:** text/plain <br />
      

Example using cURL:

```
curl --request POST \
--url http://localhost:4002/api/image \
--header 'Content-Type: multipart/form-data;' \
--form uploadfile=@/downloads/golang.png
```



### Get an image


Get an image previously uploaded. The user should use the identifier given when the upload finished.


* **URL**

  /api/image/{id}

* **Method:**

  `GET`

* **URL Params**

  ```{id}```: Id of the image to retrieve.

* **Data Params**

  None

* **Success Response:**

    * **Code:** 200 <br />
      **Content-Type:**	image/* <br />
      **Content:** an array of bytes
    


* **Error Response:**

    * **Code:** 404 NOT FOUND <br />
    **Content-Type:** application/problem+json <br />
      **Content:** <br /> 
      ```
      {
        "title": "image not found",
        "status": 404,
        "detail": "error getting image with id 974c2123adfc79c45b4dad7fc2740f4c1bbea0e90ff - image not found in storage",
        "instance": "/api/image/974c2123adfc79c45b4dad7fc2740f4c1bbea0e90ff"
      }
      ```

    * **Code:** 500 SERVER ERROR <br />
      **Content-Type:** text/plain <br />

Example using cURL:

```
curl --request GET \
  --url http://localhost:4002/api/image/23d56535b3ae661fd530f5a594117d57b04280c3ceb7144bffb9f9813e89ca5d
```



## How to run the project:

Using the makefile:

``` 
make build
make run 
```

Using the Dockerfile:

``` docker build --tag hipeople_tsk . && docker run -p 4002:4002 hipeople_tsk ```

Running tests:

``` 
make build
make test 
```
