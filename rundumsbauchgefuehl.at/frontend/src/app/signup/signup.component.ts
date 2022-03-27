import { Component, OnInit, ViewChild } from '@angular/core';
import { FormControl, Validators, FormGroup, FormGroupDirective, NgForm} from '@angular/forms';
import { SignUp } from '../sign-up.model'
import { DataService } from '../data.service'
import { MatDialog } from '@angular/material/dialog';
import { ReCaptchaV3Service } from 'ng-recaptcha';
import { ErrorDialogComponent } from '../error-dialog/error-dialog.component';

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.scss']
})

export class SignupComponent implements OnInit {

  form!: FormGroup;
  mainObject: SignUp = <SignUp>{}
  apiURL = "https://invalid";
  thankYou = "";


  constructor(
    private recaptchaV3Service: ReCaptchaV3Service,
    public dataService: DataService,
    public dialog: MatDialog
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
      pickup: new FormControl(this.mainObject.pickup),
      couponCode: new FormControl(this.mainObject.couponCode),
    });
    console.log(this.form)

  }

  @ViewChild(FormGroupDirective) formDirective!: FormGroupDirective;
  submitForm(): void {
    console.log(this.form.value)
    this.recaptchaV3Service.execute('verify_bauchgefuehl')
    .subscribe(
      (token: string) => {
        console.log(`Token [${token}] generated`);
        let v: SignUp = this.form.value
        v.reCaptchaToken = token
        this.dataService.sendData(v)
          .subscribe(
            (stream: any) => {
              console.log('send successful')

              this.form.reset();
              this.thankYou = "Danke für Ihre Bestellung! Wir haben Ihnen eine E-Mail als Bestätigung geschickt!"
              this.formDirective.resetForm();
              this.buildForm()
            },
            (error: any)  => {
              console.error(error)
              this.dialog.open(ErrorDialogComponent, {
                data: { error: error },
              });
            }
          );
         
        },
        error => {
          console.log(error);
        });
  }
}
