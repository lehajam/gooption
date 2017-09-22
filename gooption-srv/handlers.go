

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/protobuf/proto"
	"github.com/gooption/gobs"
	"github.com/gooption/pb"
)


func handlerPrice(c *gin.Context) {
	//MIMEPROTOBUF || MIMEJSON
  	request := &pb.PriceRequest{}
	if err := c.Bind(request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	response, err := gooption.Service{}.Price(request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if c.ContentType() == binding.MIMEJSON {
		c.JSON(http.StatusOK, response)
	} else {
		stream, err := proto.Marshal(response)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, binding.MIMEPROTOBUF, stream)
	}
}

func handlerGreek(c *gin.Context) {
	//MIMEPROTOBUF || MIMEJSON
  	request := &pb.GreekRequest{}
	if err := c.Bind(request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	response, err := gooption.Service{}.Greek(request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if c.ContentType() == binding.MIMEJSON {
		c.JSON(http.StatusOK, response)
	} else {
		stream, err := proto.Marshal(response)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, binding.MIMEPROTOBUF, stream)
	}
}

func handlerImpliedVol(c *gin.Context) {
	//MIMEPROTOBUF || MIMEJSON
  	request := &pb.ImpliedVolRequest{}
	if err := c.Bind(request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	response, err := gooption.Service{}.ImpliedVol(request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if c.ContentType() == binding.MIMEJSON {
		c.JSON(http.StatusOK, response)
	} else {
		stream, err := proto.Marshal(response)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, binding.MIMEPROTOBUF, stream)
	}
}
