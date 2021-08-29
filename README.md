# HiPeople Technical Task

A REST API that allows user to upload images and access them later.



## API routes

Documentation on the entry points of this web API.

### Upload an image

---

Upload an image. An identifier is returned to be used later on when retrieving the image. 
Size of image is limited to 1MB.


* **URL**

  /api/image

* **Method:**

  `POST`

* **URL Params**

  None

* **Data Params**

  None

* **Success Response:**

    * **Code:** 200 <br />
      **Content:** `{ id : 12 }`

* **Error Response:**

    * **Code:** 404 NOT FOUND <br />
      **Content:** `{ error : "User doesn't exist" }`


### Get an image

---

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
      **Content:** `{ error : "User doesn't exist" }`