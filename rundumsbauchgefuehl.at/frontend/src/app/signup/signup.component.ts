import { Component, OnInit } from '@angular/core';
import { FormControl, Validators, FormGroup } from '@angular/forms';
import { SignUp } from '../sign-up.model'
import { MatDialog } from '@angular/material/dialog';
// import { DataService } from '../data.service';

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.scss']
})

export class SignupComponent implements OnInit {

  form!: FormGroup;
  // formDirective: FormGroupDirective = <FormGroupDirective>{};
  mainObject: SignUp = <SignUp>{}
  apiURL = "https://invalid";
  thankYou = "";


  constructor(
    // public dataService: DataService,
    // public dialog: MatDialog
  ) { }


  ngOnInit(): void {
    this.buildForm();
  }

  buildForm() {
    this.form = new FormGroup({
      name: new FormControl(this.mainObject.name, Validators.required),
      address: new FormControl(this.mainObject.address, Validators.required),
      emailAddress: new FormControl(this.mainObject.emailAddress, [Validators.required, Validators.email]),
      comment: new FormControl(this.mainObject.comment),
    });

  }

  submitForm(): void {
    console.log(this.form.value)
    // TODO: Backend call!

    /*
    this.dialog.open(ErrorDialogComponent, {
      data: { error: error },
    });*/
    // this.mainObject = <SignUp>{}
    this.form.reset();
    this.thankYou = "Danke für Ihre Bestellung!"
    this.buildForm()
    /*
    this.dataService.createObject(this.apiURL, o)
      .subscribe(
        (stream: any) => {
          console.log('create successful')
          this.created(stream);
        },
        (error: any)  => {
          console.error(error)
          this.dialog.open(ErrorDialogComponent, {
            data: { error: error },
          });
        }
      );
     */

  }


}
