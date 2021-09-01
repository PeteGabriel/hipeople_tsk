# HiPeople Technical Task

A web API that allows users to upload images and access them later. 



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

  * **Multipart Form:** uploadfile <br />
    **Content:** a binary file (up to 1MB)

* **Headers**

  * **Content-Type:** 'multipart/form-data;'

* **Success Response:**

    * **Code:** 201 <br />
      **Content-Type:** text/plain; charset=utf-8 <br />
      **Content:** `{
      "image_id": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
      }`

* **Error Response:**

    * **Code:** 404 NOT FOUND <br />
      **Content:** `{ error : "User doesn't exist" }`


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

  None

* **Data Params**

  None

* **Success Response:**

    * **Code:** 200 <br />
      **Content:** `{ id : 12 }`

* **Error Response:**

    * **Code:** 404 NOT FOUND <br />
    **Content-Type:** application/problem+json <br />
      **Content:** <br /> 
      ```
      {
        "title": "image not found",
        "status": 404,
        "detail": "error getting image with id 213124125235325 - image not found",
        "instance": "/api/image/213124125235325"
      }
      ```


Example using cURL:

```
curl --request GET \
  --url http://localhost:4002/api/image/23d56535b3ae661fd530f5a594117d57b04280c3ceb7144bffb9f9813e89ca5d
```



## how to run the project:

Using the makefile:

``` make run ```

Using the Dockerfile:

``` docker build --tag hipeople_tsk . && docker run -p 4002:4002 hipeople_tsk ```