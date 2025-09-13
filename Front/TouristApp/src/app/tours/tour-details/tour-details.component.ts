import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import * as L from 'leaflet';
import 'leaflet-routing-machine'; // for road routing
import { Tour } from '../model/tour.model';
import { Keypoint } from '../model/keypoint.model';
import { ToursService } from '../tours.service';
import { forkJoin, Observable } from 'rxjs';

@Component({
  selector: 'app-tour-details',
  templateUrl: './tour-details.component.html',
  styleUrls: ['./tour-details.component.css']
})
export class TourDetailsComponent implements OnInit {
  tour!: Tour;
  keypoints: Keypoint[] = [];
  map!: L.Map;
  @ViewChild('mapContainer') mapContainer!: ElementRef;
  routeControl!: any;
  keypointLayer!: L.LayerGroup

  constructor(
    private route: ActivatedRoute,
    private tourService: ToursService
  ) {}

  ngOnInit(): void {
    const tourId = this.route.snapshot.paramMap.get('id');
    if (tourId) this.loadTour(tourId);
  }

  async loadTour(tourId: string) {
    try {
      const tourResponse = await this.tourService.getTourByID(tourId).toPromise();
      if(tourResponse){
        this.tour = tourResponse;
      }
      const keypointsResponse = await this.tourService.getKeyPointsForTour(tourId).toPromise();
      if(keypointsResponse){
        this.keypoints = keypointsResponse;
      }
      this.initMap();
    } catch (error) {
      console.error('Failed to load tour:', error);
    }
  }

  initMap() {
    if (!this.mapContainer) return;
    console.log("mapa tura")
    this.map = L.map(this.mapContainer.nativeElement).setView([44.817, 20.456], 12);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
      
    }).addTo(this.map);

    this.drawKeypoints();
    this.drawRoute();
    this.enableAddingKeypoints();
  }

  drawKeypoints() {
    // Create layer group if not exists
  if (!this.keypointLayer) {
    this.keypointLayer = L.layerGroup().addTo(this.map);
  } else {
    this.keypointLayer.clearLayers(); // remove old markers
  }

  this.keypoints.forEach(kp => {
    const bigIcon = L.icon({
      iconUrl: 'assets/markinjo.png',   // use default or custom
      iconSize: [40, 65],                  // bigger marker
      iconAnchor: [20, 65],                // so it "sits" on point
      popupAnchor: [0, -60],
    });
    const marker = L.marker([kp.latitude, kp.longitude], {
      draggable: true,
      icon: bigIcon,
      riseOnHover: true // makes marker always stay "on top"
    });//.addTo(this.map);
    this.keypointLayer.addLayer(marker); // ✅ only add to layer group
    // HTML form inside popup
    const popupContent = `
      <div style="min-width:200px">
        <h4>Edit Keypoint</h4>
        <label>Name:</label><br>
        <input id="kp-name-${kp.id}" type="text" value="${kp.name}" style="width:100%"/><br>
        
        <label>Description:</label><br>
        <textarea id="kp-desc-${kp.id}" style="width:100%">${kp.description || ''}</textarea><br>
        
        <label>Image:</label><br>
        <input id="kp-img-${kp.id}" type="file" accept="image/*" /><br><br>
        
        <button id="kp-save-${kp.id}" style="margin-top:5px">💾 Save</button>
      </div>
    `;

    marker.bindPopup(popupContent);

    marker.on('popupopen', () => {
      const saveBtn = document.getElementById(`kp-save-${kp.id}`);
      if (saveBtn) {
        saveBtn.addEventListener('click', () => {
          const nameInput = (document.getElementById(`kp-name-${kp.id}`) as HTMLInputElement).value;
          const descInput = (document.getElementById(`kp-desc-${kp.id}`) as HTMLTextAreaElement).value;
          const imgInput = (document.getElementById(`kp-img-${kp.id}`) as HTMLInputElement).files?.[0];

          kp.name = nameInput;
          kp.description = descInput;

          // build form data for backend (if image upload is needed)
          const formData = new FormData();
          formData.append("name", kp.name);
          formData.append("description", kp.description || "");
          formData.append("latitude", kp.latitude.toString());
          formData.append("longitude", kp.longitude.toString());
          if (imgInput) formData.append("image", imgInput);

          this.tourService.updateKeypoint(kp.id!, formData).subscribe(() => {
            marker.closePopup();
          });
        });
      }
    });

    marker.on('dragend', (event: any) => this.updateKeypointPosition(kp, event));
    kp['marker'] = marker;
  });

  // ✅ clear separation: add new keypoint only when clicking on the map background
  this.map.on('click', (event: L.LeafletMouseEvent) => {
    const latlng = event.latlng;
    this.createKeypoint(latlng.lat, latlng.lng);
  });
}

createKeypoint(lat: number, lng: number) {
  const newKeypoint: Keypoint = {
    id: undefined, // backend will assign
    name: "",
    description: "",
    latitude: lat,
    longitude: lng,
    images: []
  };

  const marker = L.marker([lat, lng], { draggable: true }).addTo(this.map);
  const popupContent = `
    <div style="min-width:200px">
      <h4>New Keypoint</h4>
      <label>Name:</label><br>
      <input id="new-kp-name" type="text" style="width:100%"/><br>
      
      <label>Description:</label><br>
      <textarea id="new-kp-desc" style="width:100%"></textarea><br>
      
      <label>Image:</label><br>
      <input id="new-kp-img" type="file" accept="image/*" /><br><br>
      
      <button id="new-kp-save">➕ Add</button>
    </div>
  `;

  marker.bindPopup(popupContent).openPopup();

  marker.on('popupopen', () => {
    const saveBtn = document.getElementById("new-kp-save");
    if (saveBtn) {
      saveBtn.addEventListener("click", () => {
        const nameInput = (document.getElementById("new-kp-name") as HTMLInputElement).value;
        const descInput = (document.getElementById("new-kp-desc") as HTMLTextAreaElement).value;
        const imgInput = (document.getElementById("new-kp-img") as HTMLInputElement).files?.[0];

        newKeypoint.name = nameInput;
        newKeypoint.description = descInput;

        const formData = new FormData();
        formData.append("name", newKeypoint.name);
        formData.append("description", newKeypoint.description || "");
        formData.append("latitude", newKeypoint.latitude.toString());
        formData.append("longitude", newKeypoint.longitude.toString());
        if (imgInput) formData.append("image", imgInput);

        this.tourService.createKeypoint(formData).subscribe(savedKp => {
          newKeypoint.id = savedKp.id;
          this.keypoints.push(newKeypoint);
          marker.closePopup();
        });
      });
    }
  });

  marker.on("dragend", (event: any) => this.updateKeypointPosition(newKeypoint, event));
  newKeypoint['marker'] = marker;
}


  drawRoute() {
    if (this.routeControl) this.routeControl.remove();

    const waypoints = this.keypoints!.map(kp => L.latLng(kp.latitude, kp.longitude));
    if (waypoints.length < 2) return;

    this.routeControl = L.Routing.control({
      waypoints,
      lineOptions: {
        styles: [{ color: 'blue', weight: 4 }],
        extendToWaypoints: false,
        missingRouteTolerance: 0
      },
      router: L.Routing.osrmv1({ serviceUrl: 'https://router.project-osrm.org/route/v1' }),
      //draggableWaypoints: false,
      addWaypoints: false
    }).addTo(this.map);
  }

  enableAddingKeypoints() {
    this.map.on('click', (e: L.LeafletMouseEvent) => {
      const newKp: Keypoint = {
        name: 'New Keypoint',
        description: '',
        latitude: e.latlng.lat,
        longitude: e.latlng.lng,
        images: [],
        order: this.keypoints!.length + 1
      };
      this.keypoints!.push(newKp);
      this.drawKeypoints();
      this.drawRoute();
    });
  }

  updateKeypointPosition(kp: Keypoint, event: any) {
    const latlng = event.target.getLatLng();
    kp.latitude = latlng.lat;
    kp.longitude = latlng.lng;
    this.drawRoute();
  }

  deleteKeypoint(kp: Keypoint) {
    if (kp['marker']) this.map.removeLayer(kp['marker']);
    this.keypoints = this.keypoints!.filter(k => k !== kp);
    this.drawRoute();
  }

  saveTourUpdates() {
  if (!this.tour) return;

  const updateCalls: Observable<any>[] = [];

  this.keypoints.forEach((kp, index) => {
    const formData = new FormData();
    formData.append("name", kp.name);
    formData.append("description", kp.description || "");
    formData.append("latitude", kp.latitude.toString());
    formData.append("longitude", kp.longitude.toString());
    formData.append("order", index.toString()); // 👈 maintain order
    formData.append("tourId", this.tour.id!);

    if (kp.images) {
      formData.append("images", String(kp.images));
    }

    if (kp.id) {
      // Existing → update
      updateCalls.push(this.tourService.updateKeypoint(kp.id, formData));
    } else {
      // New → create
      updateCalls.push(this.tourService.createKeypoint(formData));
    }
  });

  // Wait for all updates/creates
  forkJoin(updateCalls).subscribe({
    next: (results) => {
      console.log("All keypoints saved in order:", results);

      // After all keypoints persisted → update tour length
      this.tourService.updateTourLength(this.tour!.id!).subscribe({
        next: () => console.log("Tour length updated"),
        error: (err) => console.error("Tour length update failed", err)
      });
    },
    error: (err) => console.error("Keypoint save failed", err)
  });
}

}

