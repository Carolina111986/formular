import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { SignUp } from './sign-up.model'
import { environment } from './../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class DataService {

  private API_SERVER = environment.apiUrl;

  constructor(
    private http: HttpClient
  ) { }

  public sendData(o: SignUp) {
    let server = this.API_SERVER;
    return this.http.post(server, o); 
  }


}
