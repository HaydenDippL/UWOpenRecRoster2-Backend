import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class SchedulesService {
  private scheduleData?: Schedules;
  private scheduleDate?: string;
  private scheduleFacilities?: string;

  constructor(private http: HttpClient) {}

  fetchSchedule(date: string): {

  }

  setFacility(facility: string): void {

  }

  getSchedule(gym: string): Observable<Event[]> {

  }
}

export interface Schedules {
  bakke: FacilityEvents;
  nick: FacilityEvents;
}

export interface FacilityEvents {
  courts: Event[];
  pool: Event[];
  esports: Event[];
  mt_mendota: Event[];
  ice_rink: Event[];
}

export interface Event {
  name: string;
  location: string;
  start: string;
  end: string;
}
