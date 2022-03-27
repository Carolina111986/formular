import { Component, OnInit, Inject } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatDialogRef } from '@angular/material/dialog';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'app-error-dialog',
  templateUrl: './error-dialog.component.html',
  styleUrls: ['./error-dialog.component.scss']
})
export class ErrorDialogComponent implements OnInit {

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: {error: HttpErrorResponse}
  ) { }


  ngOnInit(): void {
  }

  errorText(): string {
    if ( 500 <= this.data.error.status && this.data.error.status < 600) {
      return "Server Fehler. Fehlercode " + this.data.error.status + ". Bitte probieren Sie es später nochmals oder bestellen Sie direkt per E-Mail bei office@rundumsbauchgefuehl.at.";
    } else if (403 === this.data.error.status) {
      return "Anmeldefehler. Bitte aus- und wieder einloggen.";
    } else if (400 <= this.data.error.status && this.data.error.status < 500) {
      return "Fehlercode " + this.data.error.status + ".";
    } else {
      return "Fehlercode " + this.data.error.status + ".";
    }
  }

}
