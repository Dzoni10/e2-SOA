import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import * as L from 'leaflet';
import { interval, Subscription } from 'rxjs';
import { ToursService } from '../tours.service';
import { Keypoint } from '../model/keypoint.model';
import { TourExecution } from '../model/tour-execution';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatProgressBar } from '@angular/material/progress-bar';
import { Tour } from '../model/tour.model';



const iconUrl = 'assets/location.png';
const iconDefault = L.icon({
  iconUrl,
  iconSize: [40, 40],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  tooltipAnchor: [16, -28]
});
L.Marker.prototype.options.icon = iconDefault;


@Component({
  selector: 'app-tour-execution',
  templateUrl: './tour-execution.component.html',
  styleUrls: ['./tour-execution.component.css']
})
export class TourExecutionComponent implements OnInit,OnDestroy {

  executionId!: string;
  execution?: TourExecution;
  map!: L.Map;
  marker!: L.Marker;
  keyPoints: Keypoint[] = [];
  checkSub?: Subscription;
  tour?: Tour; 
  keyPointMarkers: Map<string,L.Marker>=new Map();


  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private service: ToursService,
    private snackBar :MatSnackBar
  ){}

  ngOnInit(): void {
    
    this.executionId = this.route.snapshot.paramMap.get('id')!;
    this.loadExecution();

    this.map = L.map('map').setView([44.8176, 20.4569], 14);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(this.map);

    // marker turiste
    this.marker = L.marker([44.8176, 20.4569]).addTo(this.map);

    // simulacija: proverava na svakih 10 sekundi
    this.checkSub = interval(10000).subscribe(() => {
      //this.updatePosition();
    });

    this.map.on('click',(e: L.LeafletMouseEvent)=>{
      const lat = e.latlng.lat
      const lon = e.latlng.lng

      this.marker.setLatLng([lat,lon]);


      this.service.updateLocation(this.executionId,lat,lon).subscribe(()=>{
        this.snackBar.open("Location updated","Close",{duration:3000,horizontalPosition:'center'});

        this.keyPoints.forEach(kp=>{
          const dist = this.distance(lat,lon,kp.latitude,kp.longitude);
          if(dist<30 && !this.execution?.completedKeyPoints.some(c=>c.keyPointId === kp.id)){
            this.service.addCompletedKeyPoint(this.executionId,kp.id!,this.keyPoints.length).subscribe(()=>{
              this.execution?.completedKeyPoints.push({ keyPointId: kp.id!,reachedAt: new Date().toISOString() });
              this.markKeyPointCompleted(kp);
              this.updateProgressBar();
              this.snackBar.open(`${kp.name} completed!`, "Close", { duration: 3000, horizontalPosition: 'center' });
            });
          }
        });
      });

    });
  }

  ngOnDestroy(): void {
    this.checkSub?.unsubscribe();
  }

  loadExecution(): void {
    this.service.getExecutionById(this.executionId).subscribe(ex => {
      this.execution = ex;

      this.marker = L.marker([ex.currentLatitude,ex.currentLongitude], {icon: iconDefault})
      .addTo(this.map)
      .bindPopup('Your Location')
      .openPopup();

      this.service.getTourByID(ex.tourId).subscribe(t => {
        this.tour = t;
      });
      this.loadKeyPoints(ex.tourId);

      this.updateProgressBar();
    });
  }
  
  loadKeyPoints(tourId: string): void {
    // pozoveš svoj KeyPointService
    // dummy primer:
    // this.keyPoints = [
    //   { id: '1', name: 'Gate', description: '', latitude: 44.817, longitude: 20.457,images:[] },
    //   { id: '2', name: 'Museum', description: '', latitude: 44.818, longitude: 20.455,images:[] }
    // ];

    this.service.getKeyPointsForTour(tourId).subscribe(kps => {
    this.keyPoints = kps;

    this.keyPoints.forEach(kp => {
      const marker = L.marker([kp.latitude, kp.longitude], {
        icon: L.icon({ iconUrl: 'assets/markinjo.png', iconSize: [30, 30] })
      }).addTo(this.map)
        .bindPopup(kp.name);

      this.keyPointMarkers.set(kp.id!, marker);

      // ako je već pređen
      if (this.execution?.completedKeyPoints.some(c => c.keyPointId === kp.id)) {
        this.markKeyPointCompleted(kp);
      }
    });
  }); 
  }

  markKeyPointCompleted(kp: Keypoint) {
  const marker = this.keyPointMarkers.get(kp.id!);
  if (marker) {
    marker.setIcon(L.icon({ iconUrl: 'assets/markinjo_done.png', iconSize: [24, 24] }));
  }
}

updateProgressBar() {
  if (!this.execution) return;
  const completed = this.execution.completedKeyPoints.length;
  const total = this.keyPoints.length;
  this.execution.completedPercent =Math.min(100,(completed / total) * 100);
}

  updatePosition(): void {
    // dobavi novu lokaciju od PositionSimulatora
    const lat = 44.817 + (Math.random() - 0.5) * 0.002;
    const lon = 20.456 + (Math.random() - 0.5) * 0.002;

    this.marker.setLatLng([lat, lon]);

    this.service.updateLocation(this.executionId, lat, lon).subscribe();

    // proveri blizinu keypointa
    this.keyPoints.forEach(kp => {
      const dist = this.distance(lat, lon, kp.latitude, kp.longitude);
      if (dist < 30 && !this.execution?.completedKeyPoints.some(c=>c.keyPointId === kp.id)) { // npr. 30 metara prag
        this.service.addCompletedKeyPoint(this.executionId, kp.id!, this.keyPoints.length).subscribe(() => {
          this.execution?.completedKeyPoints.push({
            keyPointId:kp.id!,
            reachedAt: new Date().toISOString()
          });
          this.markKeyPointCompleted(kp)
          this.updateProgressBar();
        });
      }
    });
  }

  distance(lat1: number, lon1: number, lat2: number, lon2: number): number {
    const R = 6371e3; // metres
    const φ1 = lat1 * Math.PI/180;
    const φ2 = lat2 * Math.PI/180;
    const Δφ = (lat2-lat1) * Math.PI/180;
    const Δλ = (lon2-lon1) * Math.PI/180;

    const a = Math.sin(Δφ/2) * Math.sin(Δφ/2) +
              Math.cos(φ1) * Math.cos(φ2) *
              Math.sin(Δλ/2) * Math.sin(Δλ/2);
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));

    return R * c; // in metres
  }

  finish(): void {
    this.service.finishTour(this.executionId).subscribe(() => {
      this.snackBar.open("Tour Finished","Close",{duration:3000,horizontalPosition:"center"})
      this.router.navigate(['/allTours']);
    });
  }

  abandon(): void {
    this.service.abandonTour(this.executionId).subscribe(() => {
      this.snackBar.open("Tour abandoned","Close",{duration:3000,horizontalPosition:"center"})
      this.router.navigate(['/allTours']);
    });
  }

  get canFinish():boolean{
    return this.execution?.completedPercent ===100;
  }

}
