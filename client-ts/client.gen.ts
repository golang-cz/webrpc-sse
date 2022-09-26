/* eslint-disable */
// webrpc-sse v1.0.0 a53b3b444366958a0362df8e6fe2a251ac51cb18
// --
// This file has been generated by https://github.com/webrpc/webrpc using gen/typescript
// Do not edit by hand. Update your webrpc schema and re-generate.

// WebRPC description and code-gen version
export const WebRPCVersion = "v1"

// Schema version of your RIDL schema
export const WebRPCSchemaVersion = "v1.0.0"

// Schema hash generated from your RIDL schema
export const WebRPCSchemaHash = "a53b3b444366958a0362df8e6fe2a251ac51cb18"


//
// Types
//
export interface Message {
  id: number
  msg: string
  author: string
  created_at: string
}

export interface Chatbot {
  sendMessage(args: SendMessageArgs, headers?: object): Promise<SendMessageReturn>
  subscribeMessages(headers?: object): Promise<SubscribeMessagesReturn>
}

export interface SendMessageArgs {
  author: string
  msg: string
}

export interface SendMessageReturn {
  success: boolean  
}
export interface SubscribeMessagesArgs {
}

export interface SubscribeMessagesReturn {
  msgs: Array<Message>  
}


  
//
// Client
//
export class Chatbot implements Chatbot {
  protected hostname: string
  protected fetch: Fetch
  protected path = '/rpc/Chatbot/'

  constructor(hostname: string, fetch: Fetch) {
    this.hostname = hostname
    this.fetch = fetch
  }

  private url(name: string): string {
    return this.hostname + this.path + name
  }
  
  sendMessage = (args: SendMessageArgs, headers?: object): Promise<SendMessageReturn> => {
    return this.fetch(
      this.url('SendMessage'),
      createHTTPRequest(args, headers)).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          success: <boolean>(_data.success)
        }
      })
    })
  }
  
  subscribeMessages = (headers?: object): Promise<SubscribeMessagesReturn> => {
    return this.fetch(
      this.url('SubscribeMessages'),
      createHTTPRequest({}, headers)
      ).then((res) => {
      return buildResponse(res).then(_data => {
        return {
          msgs: <Array<Message>>(_data.msgs)
        }
      })
    })
  }
  
}

  
export interface WebRPCError extends Error {
  code: string
  msg: string
	status: number
}

const createHTTPRequest = (body: object = {}, headers: object = {}): object => {
  return {
    method: 'POST',
    headers: { ...headers, 'Content-Type': 'application/json' },
    body: JSON.stringify(body || {})
  }
}

const buildResponse = (res: Response): Promise<any> => {
  return res.text().then(text => {
    let data
    try {
      data = JSON.parse(text)
    } catch(err) {
      throw { code: 'unknown', msg: `expecting JSON, got: ${text}`, status: res.status } as WebRPCError
    }
    if (!res.ok) {
      throw data // webrpc error response
    }
    return data
  })
}

export type Fetch = (input: RequestInfo, init?: RequestInit) => Promise<Response>