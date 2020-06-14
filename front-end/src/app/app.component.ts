import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { interval, Observable } from 'rxjs';
import { startWith, switchMap } from "rxjs/operators";
import { config } from 'process';


const server = '192.168.0.112';
const port = '8080';
const protocol = 'http';

const config = {
  prod: {
    base: `/`,
  },
  dev: {
    base: `${protocol}://${server}:${port}/`,
  }
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent {
  urlProcs = config.prod.base + 'procs';
  urlCPU = config.prod.base + 'cpu';
  urlMEMO = config.prod.base + 'memo';

  procesos: string[] = [];
  ejecucion: number = 0;
  suspendido: number = 0;
  detenido: number = 0;
  zombie: number = 0;
  total: number = 0;

  cpuUsed: number = 0;

  cpuMemoTotalMb: number = 0;
  cpuMemoUsedMb: number = 0;

  constructor(private http: HttpClient) {
    interval(15000)
      .pipe(
        startWith(0),
        switchMap(() => this.getProcs())
      )
      .subscribe(res => {
        this.suspendido = 0;
        this.procesos.length = 0;
        for (let key in res) {
          if (res.hasOwnProperty(key)) {
            this.procesos.push(res[key]);
            if (res[key].EstadoID === 'S') {
              this.suspendido++;
            } else if (res[key].EstadoID === 'I') {
            }
          }
        }
        this.total = this.procesos.length;
      });

    interval(3000)
      .pipe(
        startWith(0),
        switchMap(() => this.getCPU())
      )
      .subscribe(res => {
        for (let key in res) {
          this.cpuUsed = res[key];
        }
      });

    interval(2000)
      .pipe(
        startWith(0),
        switchMap(() => this.getMemo())
      )
      .subscribe(res => {
        this.cpuMemoTotalMb = res['TotalMb'];
        this.cpuMemoUsedMb = res['UsedMb'];
      });
  }

  getProcs(): Observable<any> {
    return this.http
      .get(this.urlProcs);
  }

  getCPU(): Observable<any> {
    return this.http
      .get(this.urlCPU);
  }

  getMemo(): Observable<any> {
    return this.http
      .get(this.urlMEMO);
  }
}
